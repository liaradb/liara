package scan

import "encoding/binary"

func Int16(data []byte) (int16, []byte, bool) {
	if len(data) < 2 {
		return 0, nil, false
	}

	return int16(binary.BigEndian.Uint16(data[:2])), data[2:], true
}

func Uint16(data []byte) (uint16, []byte, bool) {
	if len(data) < 2 {
		return 0, nil, false
	}

	return binary.BigEndian.Uint16(data[:2]), data[2:], true
}

func SetInt16(data []byte, v int16) ([]byte, bool) {
	if len(data) < 2 {
		return nil, false
	}

	binary.BigEndian.PutUint16(data[:2], uint16(v))
	return data[2:], true
}

func SetUint16(data []byte, v uint16) ([]byte, bool) {
	if len(data) < 2 {
		return nil, false
	}

	binary.BigEndian.PutUint16(data[:2], v)
	return data[2:], true
}
