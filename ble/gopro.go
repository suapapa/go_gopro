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
	"github.com/suapapa/go_gopro/open_gopro"
	"google.golang.org/protobuf/proto"
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

// KeepAlive sends a keep alive message to the GoPro.
// The best practice to prevent the GoPro from sleeping is to send a keep alive message every 3 seconds.
func (g *GoPro) KeepAlive() error {
	reqPayload := []byte{0x5b, 0x01, 0x42}
	resp, err := g.doRequest(
		Setting, SettingResponse,
		reqPayload,
		5*time.Second,
	)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}

	expectedRespPayload := []byte{0x5b, 0x00}
	if bytes.Compare(resp, expectedRespPayload) != 0 {
		return fmt.Errorf("unexpected response, %x", resp)
	}
	return nil
}

// SetShutter sets the shutter on or off.
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

	resp, err := g.doRequest(
		Command, CommandResponse,
		reqPayload,
		5*time.Second,
	)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}

	expectedRespPayload := makeTlvResp(cmdSetShutter, cmdRespSuccess, nil)
	if bytes.Compare(resp, expectedRespPayload) != 0 {
		return fmt.Errorf("unexpected response, %x", resp)
	}

	return nil
}

// Sleep puts the camera into sleep mode.
func (g *GoPro) Sleep() error {
	reqPayload := makeTlvCmd(cmdSleep)

	resp, err := g.doRequest(
		Command, CommandResponse,
		reqPayload,
		5*time.Second,
	)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}

	expectedRespPayload := makeTlvResp(cmdSleep, cmdRespSuccess, nil)
	if bytes.Compare(resp, expectedRespPayload) != 0 {
		return fmt.Errorf("unexpected response, %x", resp)
	}

	return nil
}

// SetTime sets the date and time on the camera.
func (g *GoPro) SetTime(t time.Time) error {
	reqPayload, err := makeTlvCmdWithParam(cmdSetDateTime, time2Bytes(t))
	if err != nil {
		return errors.Wrap(err, "failed to make tlv")
	}

	resp, err := g.doRequest(
		Command, CommandResponse,
		reqPayload,
		5*time.Second,
	)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}

	expectedRespPayload := makeTlvResp(cmdSetDateTime, cmdRespSuccess, nil)
	if bytes.Compare(resp, expectedRespPayload) != 0 {
		return fmt.Errorf("unexpected response, %x", resp)
	}

	return nil
}

// GetTime gets the date and time on the camera.
func (g *GoPro) GetTime() (time.Time, error) {
	reqPayload := makeTlvCmd(cmdGetDateTime)

	resp, err := g.doRequest(
		Command, CommandResponse,
		reqPayload,
		5*time.Second,
	)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to request")
	}

	if len(resp) != 11 {
		return time.Time{}, fmt.Errorf("unexpected response, %x", resp)
	}

	if resp[0] != cmdGetDateTime || resp[1] != cmdRespSuccess {
		return time.Time{}, fmt.Errorf("unexpected response, %x", resp)
	}

	// parsing resp.
	t, err := bytes2Time(resp[3:])
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to parse time")
	}

	return t, nil
}

/*
func (g *GoPro) SetLocalTime(t time.Time, loc time.Location) error {
	// TBU
}

func (g *GoPro) GetLocalTime() (time.Time, error) {
	// TBU
}
*/

func (g *GoPro) SetLivestreamMode(mode *open_gopro.RequestSetLiveStreamMode) error {
	pbPayload, err := proto.Marshal(mode)
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}
	reqPayload := append([]byte{featureCommand, actionSetLiveStream}, pbPayload...)

	pbResp, err := g.doRequest(
		Command, CommandResponse,
		reqPayload,
		5*time.Second,
	)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}

	if pbResp[0] != featureCommand || pbResp[1] != responseSetLiveStream {
		return fmt.Errorf("unexpected response, %x", pbResp)
	}

	resp := &open_gopro.ResponseGeneric{}
	err = proto.Unmarshal(pbResp[2:], resp)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}

	if resp.GetResult() != open_gopro.EnumResultGeneric_RESULT_SUCCESS {
		return fmt.Errorf("failed to set live stream mode, %s", resp)
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
