package ble

import (
	"io"

	"github.com/pkg/errors"
)

// first byte of packet header
type packetHeaderByte1 byte

const (
	packetHeaderContinuationMask = 0b1000_0000

	blePackaeHeaderMessageTypeMask     = 0b0110_0000
	packetHeaderMessageType5bitLength  = 0b0000_0000
	packetHeaderMessageType13bitLength = 0b0010_0000
	packetHeaderMessageType16bitLength = 0b0100_0000
	packetHeaderMessageTypeReserved    = 0b0110_0000
	packetHeaderMessageLengthMask      = 0b0001_1111

	packetHeaderContinuationCounterMask = 0b0000_1111
)

// Parse parses the header byte and returns the start flag, message type and low byte
// when start is true, the low byte is varied by the message type
// when start is false, the low byte is the continuation counter
func (p packetHeaderByte1) Parse() (start bool, msgType byte, low byte) {
	if p&packetHeaderContinuationMask != 0 {
		low = byte(p & packetHeaderContinuationCounterMask)
		return
	}

	// start packet
	msgType = byte(p & blePackaeHeaderMessageTypeMask)
	switch msgType {
	case packetHeaderMessageType5bitLength:
		// low is 5bit length
		low = byte(p & packetHeaderMessageLengthMask)
	case packetHeaderMessageType13bitLength:
		// low is upper 5bit of 13bit length
		// lower 8bit of 13bit length is read from next byte
		low = byte(p & packetHeaderMessageLengthMask)
	case packetHeaderMessageType16bitLength:
		// don't use low
		// 16bit length is read from next 2 bytes
	}

	return
}

// ReadPackets reads a packets from r and returns the payload.
func ReadPackets(r io.Reader) ([]byte, error) {
	var packet []byte
	buffer := make([]byte, 20) // 20 byte is max length of BLE packet

	for {
		n, err := r.Read(buffer)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read packet")
		}
		if n == 0 {
			return nil, errors.New("failed to read packet")
		}

		start, msgType, low := packetHeaderByte1(buffer[0]).Parse()
		var msgLen int
		if start {
			switch msgType {
			case packetHeaderMessageType5bitLength:
				msgLen = int(low)
				packet = append(packet, buffer[1:n]...)
			case packetHeaderMessageType13bitLength:
				msgLen = int(low)<<8 + int(buffer[1])
				packet = append(packet, buffer[2:n]...)
			case packetHeaderMessageType16bitLength:
				msgLen = int(buffer[1])<<8 + int(buffer[2])
				packet = append(packet, buffer[3:n]...)
			}
		} else {
			if packet == nil {
				return nil, errors.New("invalid packet")
			}
			packet = append(packet, buffer[1:n]...)
		}
		if len(packet) >= msgLen {
			return packet[:msgLen], nil
		}
	}
}
