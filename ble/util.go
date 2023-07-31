package ble

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

func time2Bytes(t time.Time) []byte {
	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	return []byte{
		byte(year >> 8), byte(year),
		byte(month),
		byte(day),
		byte(h),
		byte(m),
		byte(s),
	}
}

func bytes2Time(b []byte) (time.Time, error) {
	if len(b) < 7 {
		return time.Time{}, fmt.Errorf("invalid time bytes")
	}

	year := int(b[0])<<8 | int(b[1])
	month := int(b[2])
	day := int(b[3])
	h := int(b[4])
	m := int(b[5])
	s := int(b[6])

	return time.Date(year, time.Month(month), day, h, m, s, 0, time.Local), nil
}

func uint32ToBytes(i uint32) []byte {
	return []byte{
		byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i),
	}
}

// parseBytesStr parses a string of bytes in following form into a byte array.
// 20:15:F1:79:0A:03:78:78:78:10:01:18:07:38:7B:40:95:06:48:C8:80:03:50:00
func parseBytesStr(s string) []byte {
	hex2int := func(s string) int {
		i, err := hex.DecodeString(s)
		if err != nil {
			panic(err)
		}
		return int(i[0])
	}
	var ret []byte
	for i := 0; i < len(s); i += 3 {
		ret = append(ret, byte(hex2int(s[i:i+2])))
	}
	return ret
}

func retrieveServices(a *adapter.Adapter1, dev *device.Device1) string {
	// log.Debug("Listing exposed services")
	ret := ""

	list, err := dev.GetAllServicesAndUUID()
	if err != nil {
		return ""
	}

	if len(list) == 0 {
		time.Sleep(time.Second * 2)
		return retrieveServices(a, dev)
	}

	for _, servicePath := range list {
		ret += fmt.Sprintf("- %s\n", servicePath)
	}

	return ret
}
