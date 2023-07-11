package ble

// KeepAlive sends a keep alive command to the camera
// Send KeepAlive every 3 seconds to prevent the camera from sleeping
func KeepAlive() {
	// send followin payload to GP-0074, characteristicSetting
	// p := []byte{0x03, 0x5b, 0x01, 0x42}

	// response will be sent to GP-0075, characteristicSettingResponse
	// with payload 0x02, 0x5b, 0x00
}
