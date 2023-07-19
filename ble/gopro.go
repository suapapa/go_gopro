package ble

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	reqPayload := []byte{0x5b, 0x01, 0x42}
	respPayload := []byte{0x5b, 0x00}
	resp, err := g.doRequest(
		Setting, SettingResponse,
		reqPayload,
		5*time.Second,
	)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}
	if bytes.Compare(resp, respPayload) != 0 {
		return fmt.Errorf("unexpected response, %x", resp)
	}
	return nil
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
	expectedRespPayload := makeTlvResp(cmdSetShutter, cmdRespSuccess, nil)

	resp, err := g.doRequest(
		Command, CommandResponse,
		reqPayload,
		5*time.Second,
	)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}
	if bytes.Compare(resp, expectedRespPayload) != 0 {
		return fmt.Errorf("unexpected response, %x", resp)
	}

	return nil
}

func (g *GoPro) doRequest(
	reqC, respC Characteristic,
	reqPayload []byte,
	timeout time.Duration,
) ([]byte, error) {
	chrReq, err := g.getChr(reqC)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chr")
	}
	chrResp, err := g.getChr(respC)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chr")
	}

	errCh := make(chan error)
	defer close(errCh)
	respCh := make(chan []byte)
	defer close(respCh)

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		resp, err := readPackets(pr)
		if err != nil {
			errCh <- errors.Wrap(err, "failed to read packets")
			return
		}
		respCh <- resp
	}()

	notiHandler := func(req []byte) {
		pw.Write(req)
	}
	err = g.cln.Subscribe(chrResp, false, notiHandler)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to chr")
	}
	defer g.cln.Unsubscribe(chrResp, false)

	pkts, err := makePackets(reqPayload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make packets")
	}
	for _, p := range pkts {
		err = g.cln.WriteCharacteristic(chrReq, p, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to write chr")
		}
		// time.Sleep(100 * time.Millisecond)
	}

	select {
	case err := <-errCh:
		if err != nil {
			return nil, errors.Wrap(err, "failed to receive notification")
		}
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	case resp := <-respCh:
		return resp, nil
	}
	return nil, errors.New("unreachable")
}

func (g *GoPro) getChr(c Characteristic) (*goble.Characteristic, error) {
	ch, ok := g.chs[c]
	if !ok {
		return nil, fmt.Errorf("chr %s not found", c)
	}
	return ch, nil
}
