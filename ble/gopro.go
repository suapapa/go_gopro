package ble

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"

	goble "github.com/go-ble/ble"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/pkg/errors"
	"github.com/suapapa/go_gopro/open_gopro"
	"google.golang.org/protobuf/proto"
)

type GoPro struct {
	adt *adapter.Adapter1
	dev *device.Device1

	cln goble.Client
	p   *goble.Profile
}

func ScanGoPro(adaptorID string, tmo time.Duration) ([]*GoPro, error) {
	if adaptorID == "" {
		adaptorID = adapter.GetDefaultAdapterID()
	}
	adt, err := adapter.GetAdapter(adaptorID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get adaptor")
	}

	filter := &adapter.DiscoveryFilter{
		UUIDs: []string{svcUUIDControlAndQuery},
	}
	_, cancelDiscoved, err := api.Discover(adt, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to discover")
	}
	select {
	// case dev := <-chDeviceDiscovered:
	// 	gp := &GoPro{
	// 		adt: adt,
	// 		dev: dev.Path,
	// 	}
	case <-time.After(tmo):
		cancelDiscoved()
	}

	ret := []*GoPro{}
	devs, err := adt.GetDevices()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get devices")
	}

	for _, dev := range devs {
		ret = append(ret, &GoPro{
			adt: adt,
			dev: dev,
		})
	}

	if len(ret) == 0 {
		return nil, errors.New("no GoPro found")
	}

	return ret, nil
}

func (g *GoPro) Close() error {
	g.adt.Close()

	// Unsubscribe from notifications
	exitC := g.cln.Disconnected()
	err := g.cln.CancelConnection()
	if err != nil {
		return errors.Wrap(err, "failed to cancel connection")
	}
	<-exitC

	return nil
}

func (g *GoPro) String() string {
	return fmt.Sprintf("%s: %s - %s", g.adt.Interface(), g.dev.Properties.Name, g.dev.Properties.Address)
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

// SetLivestreamMode sets the live stream mode.
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

// ApContol turn on or off AP mode
func (g *GoPro) ApControl(on bool) error {
	var param []byte
	if on {
		param = []byte{0x01}
	} else {
		param = []byte{0x00}
	}
	reqPayload, err := makeTlvCmdWithParam(cmdApControl, param)
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

	expectedRespPayload := makeTlvResp(cmdApControl, cmdRespSuccess, nil)
	if bytes.Compare(resp, expectedRespPayload) != 0 {
		return fmt.Errorf("unexpected response, %x", resp)
	}

	return nil
}

// MediaHiLightMoment highlights moment during encoding
func (g *GoPro) MediaHiLightMoment() error {
	request := makeTlvCmd(cmdMediaHiLightMoment)

	resp, err := g.doRequest(
		Command, CommandResponse,
		request,
		5*time.Second,
	)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}

	expectedRespPayload := makeTlvResp(cmdMediaHiLightMoment, cmdRespSuccess, nil)
	if bytes.Compare(resp, expectedRespPayload) != 0 {
		return fmt.Errorf("unexpected response, %x", resp)
	}

	return nil
}

type HardwareInfo struct {
	ModelID         string
	Name            string
	BoardType       string
	FirmwareVersion string
	SerialNumber    string
	ApSSID          string
	ApMAC           string
}

func parseHardwareInfo(b []byte) *HardwareInfo {
	ret := &HardwareInfo{}
	modelIDLen := int(b[0])
	ret.ModelID = string(b[1 : 1+modelIDLen])
	b = b[1+modelIDLen:]
	nameLen := int(b[0])
	ret.Name = string(b[1 : 1+nameLen])
	b = b[1+nameLen:]
	boardTypeLen := int(b[0])
	ret.BoardType = string(b[1 : 1+boardTypeLen])
	b = b[1+boardTypeLen:]
	firmwareVersionLen := int(b[0])
	ret.FirmwareVersion = string(b[1 : 1+firmwareVersionLen])
	b = b[1+firmwareVersionLen:]
	serialNumberLen := int(b[0])
	ret.SerialNumber = string(b[1 : 1+serialNumberLen])
	b = b[1+serialNumberLen:]
	apSSIDLen := int(b[0])
	ret.ApSSID = string(b[1 : 1+apSSIDLen])
	b = b[1+apSSIDLen:]
	apMACLen := int(b[0])
	ret.ApMAC = string(b[1 : 1+apMACLen])
	return ret
}

type Preset byte

const (
	PresetVideo Preset = iota
	PresetPhoto
	PresetTimelapse
)

// PresetLoadGroup loads a preset group.
func (g *GoPro) PresetLoadGroup(p Preset) error {
	var param []byte
	switch p {
	case PresetVideo:
		param = []byte{0x03, 0xE8}
	case PresetPhoto:
		param = []byte{0x03, 0xE9}
	case PresetTimelapse:
		param = []byte{0x03, 0xEA}
	default:
		return fmt.Errorf("invalid preset, %d", p)
	}
	reqPayload, err := makeTlvCmdWithParam(cmdPresetLoadGroup, param)
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

	expectedRespPayload := makeTlvResp(cmdPresetLoadGroup, cmdRespSuccess, nil)
	if bytes.Compare(resp, expectedRespPayload) != 0 {
		return fmt.Errorf("unexpected response, %x", resp)
	}

	return nil
}

// PresetLoad loads a preset with ID.
func (g *GoPro) PresetLoad(id uint32) error {
	param := uint32ToBytes(id)
	reqPayload, err := makeTlvCmdWithParam(cmdPresetLoad, param)
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

	expectedRespPayload := makeTlvResp(cmdPresetLoad, cmdRespSuccess, nil)
	if bytes.Compare(resp, expectedRespPayload) != 0 {
		return fmt.Errorf("unexpected response, %x", resp)
	}

	return nil
}

// Analytics sets third party client
func (g *GoPro) Analytics() error {
	reqPayload := makeTlvCmd(cmdAnalytics)

	resp, err := g.doRequest(
		Command, CommandResponse,
		reqPayload,
		5*time.Second,
	)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}

	expectedRespPayload := makeTlvResp(cmdAnalytics, cmdRespSuccess, nil)
	if bytes.Compare(resp, expectedRespPayload) != 0 {
		return fmt.Errorf("unexpected response, %x", resp)
	}

	return nil
}

// GetVersion returns the version of the camera.
// in form of major.minor
func (g *GoPro) GetVersion() (string, error) {
	request := makeTlvCmd(cmdOpenGoPro)

	resp, err := g.doRequest(
		Command, CommandResponse,
		request,
		5*time.Second,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to request")
	}

	if len(resp) < 3 {
		return "", fmt.Errorf("unexpected response, %x", resp)
	}

	if resp[0] != cmdOpenGoPro || resp[1] != cmdRespSuccess {
		return "", fmt.Errorf("unexpected response, %x", resp)
	}

	verStr, err := parseVersion(resp[2:])
	if err != nil {
		return "", errors.Wrap(err, "failed to parse version")
	}

	return verStr, nil
}

func parseVersion(b []byte) (string, error) {
	majorLen := int(b[0])
	major, err := strconv.Atoi(string(b[1 : 1+majorLen]))
	if err != nil {
		return "", errors.Wrap(err, "failed to parse major")
	}
	b = b[1+majorLen:]
	minorLen := int(b[0])
	minor, err := strconv.Atoi(string(b[1 : 1+minorLen]))
	if err != nil {
		return "", errors.Wrap(err, "failed to parse minor")
	}
	return fmt.Sprintf("%d.%d", major, minor), nil
}

// GetHardwareInfo gets the hardware info of the camera.
func (g *GoPro) GetHardwareInfo() (*HardwareInfo, error) {
	request := makeTlvCmd(cmdGetHardwareInfo)

	resp, err := g.doRequest(
		Command, CommandResponse,
		request,
		5*time.Second,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request")
	}

	if len(resp) < 3 {
		return nil, fmt.Errorf("unexpected response, %x", resp)
	}

	if resp[0] != cmdGetHardwareInfo || resp[1] != cmdRespSuccess {
		return nil, fmt.Errorf("unexpected response, %x", resp)
	}

	return parseHardwareInfo(resp[2:]), nil
}

func (g *GoPro) doRequest(
	reqC, respC uuid,
	reqPayload []byte,
	timeout time.Duration,
) ([]byte, error) {
	chrReq, err := g.getChrByGpUUID(reqC)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chr")
	}
	chrResp, err := g.getChrByGpUUID(respC)
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
	err = g.cln.ClearSubscriptions()
	if err != nil {
		return nil, errors.Wrap(err, "failed to clear subscriptions")
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

func (g *GoPro) getChrByGpUUID(id uuid) (*goble.Characteristic, error) {
	chr := &goble.Characteristic{
		UUID: gpUUID(id),
	}
	ret := g.p.FindCharacteristic(chr)
	if ret == nil {
		return nil, fmt.Errorf("chr %s not found", id)
	}
	return ret, nil
}
