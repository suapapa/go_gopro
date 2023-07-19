package ble

import (
	"fmt"
	"time"
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
