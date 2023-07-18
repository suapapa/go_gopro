package ble

import "fmt"

// tlv command ids
const (
	cmdSetShutter         = 0x01
	cmdSleep              = 0x05
	cmdSetDateTime        = 0x0D
	cmdGetDateTime        = 0x0E
	cmdSetLocalDataTime   = 0x0F
	cmdGetLocalDataTime   = 0x10
	cmdSetLivestreamMode  = 0x15
	cmdApControl          = 0x17
	cmdMediaHiLightMoment = 0x18
	cmdGetHardwareInfo    = 0x3C
	cmdPresetLoadGroup    = 0x3E
	cmdPresetLoad         = 0x40
	cmdAnalytics          = 0x50
	cmdOpenGoPro          = 0x51

	cmdRespSuccess          = 0
	cmdRespError            = 1
	cmdRespInvalidParameter = 2
)

func makeTlvCmd(cmd byte) []byte {
	return []byte{cmd}
}

func makeTlvCmdWithParam(cmd byte, param []byte) ([]byte, error) {
	if len(param) > 255 {
		return nil, fmt.Errorf("tlv payload too long")
	}
	return append([]byte{cmd, byte(len(param))}, param...), nil
}

func makeTlvResp(cmd byte, respCode byte, resp []byte) []byte {
	ret := []byte{cmd, respCode}
	if resp != nil {
		ret = append(ret, resp...)
	}
	return ret
}
