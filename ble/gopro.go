package ble

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	goble "github.com/go-ble/ble"
	"github.com/pkg/errors"
)

type GoPro struct {
	cln goble.Client
	p   *goble.Profile
	chs map[Characteristic]*goble.Characteristic
}

func ScanGoPro(opts ...goble.Option) (*GoPro, error) {
	dev, err := newDevice(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new device")
	}

	goble.SetDefaultDevice(dev)

	ctx := goble.WithSigHandler(context.WithTimeout(context.Background(), 10*time.Second))
	filter := func(a goble.Advertisement) bool {
		svcs := a.Services()
		for _, svc := range svcs {
			log.Println(svc)
			if svc.Equal(svcUUIDControlAndQuery) {
				return true
			}
		}
		return false
	}

	cln, err := goble.Connect(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to GoPro")
	}

	p, err := cln.DiscoverProfile(true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to discover profile")
	}

	ret := &GoPro{
		cln: cln,
		p:   p,
		chs: makeCharacteristicMap(p),
	}

	return ret, nil
}

func (g *GoPro) Close() error {
	// Unsubscribe from notifications
	exitC := g.cln.Disconnected()
	err := g.cln.CancelConnection()
	if err != nil {
		return errors.Wrap(err, "failed to cancel connection")
	}
	<-exitC

	return nil
}

func (g *GoPro) KeepAlive() error {
	return g.writePayload(
		Setting, SettingResponse,
		[]byte{0x5b, 0x01, 0x42}, []byte{0x5b, 0x00},
		5*time.Second,
	)
}

func (g *GoPro) SetShutter(on bool) error {
	var param []byte
	if on {
		param = []byte{0x01}
	} else {
		param = []byte{0x00}
	}
	reqPayload, err := makeTlvCmdWithParam(cmdSetShutter, param)
	if err != nil {
		return errors.Wrap(err, "failed to make tlv")
	}
	// TODO: make tlv resp with helper
	// expectedRespPayload := []byte{0x02, 0x01, 0x00}
	expectedRespPayload := makeTlvResp(cmdSetShutter, cmdRespSuccess, nil)

	return g.writePayload(
		Command, CommandResponse,
		reqPayload, expectedRespPayload,
		5*time.Second,
	)
}

func (g *GoPro) writePayload(
	reqC, respC Characteristic,
	reqPayload, respPayload []byte,
	timeout time.Duration,
) error {
	chrReq, err := g.getChr(reqC)
	if err != nil {
		return errors.Wrap(err, "failed to get chr")
	}
	chrResp, err := g.getChr(respC)
	if err != nil {
		return errors.Wrap(err, "failed to get chr")
	}

	doneCh := make(chan error)
	// TODO: handle multiple responses
	respPayload = append([]byte{byte(len(respPayload))}, respPayload...)

	notiHandler := func(req []byte) {
		if bytes.Equal(req, respPayload) {
			doneCh <- nil
			return
		}
		doneCh <- errors.Errorf("unexpected response: %v", req)
	}
	err = g.cln.Subscribe(chrResp, false, notiHandler)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to chr")
	}
	defer g.cln.Unsubscribe(chrResp, false)

	pkts, err := makePackets(reqPayload)
	if err != nil {
		return errors.Wrap(err, "failed to make packets")
	}
	for _, p := range pkts {
		err = g.cln.WriteCharacteristic(chrReq, p, false)
		if err != nil {
			return errors.Wrap(err, "failed to write chr")
		}
		// time.Sleep(100 * time.Millisecond)
	}

	select {
	case err := <-doneCh:
		if err != nil {
			return errors.Wrap(err, "failed to receive notification")
		}
	case <-time.After(timeout):
		return errors.New("timeout")
	}
	return nil
}

func (g *GoPro) getChr(c Characteristic) (*goble.Characteristic, error) {
	ch, ok := g.chs[c]
	if !ok {
		return nil, fmt.Errorf("chr %s not found", c)
	}
	return ch, nil
}
