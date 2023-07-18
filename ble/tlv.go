package ble

import "fmt"

// tlv command ids
const (
	commandSetShutter         = 0x01
	commandSleep              = 0x05
	commandSetDateTime        = 0x0D
	commandGetDateTime        = 0x0E
	commandSetLocalDataTime   = 0x0F
	commandGetLocalDataTime   = 0x10
	commandSetLivestreamMode  = 0x15
	commandApControl          = 0x17
	commandMediaHiLightMoment = 0x18
	commandGetHardwareInfo    = 0x3C
	commandPresetLoadGroup    = 0x3E
	commandPresetLoad         = 0x40
	commandAnalytics          = 0x50
	commandOpenGoPro          = 0x51

	commandRespSuccess          = 0
	commandRespError            = 1
	commandRespInvalidParameter = 2
)

func makeTlvCommand(cmd byte) []byte {
	return []byte{cmd}
}

func makeTlvCommandWithParam(cmd byte, param []byte) ([]byte, error) {
	if len(param) > 255 {
		return nil, fmt.Errorf("tlv payload too long")
	}
	return append([]byte{cmd, byte(len(param))}, param...), nil
}
