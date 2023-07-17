package ble

import "fmt"

// tlv command id
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

func makeTlvCommand(cmdID byte, payload []byte) ([]byte, error) {
	if len(payload) > 255 {
		return nil, fmt.Errorf("tlv payload too long")
	}
	return append([]byte{cmdID, byte(len(payload))}, payload...), nil
}
