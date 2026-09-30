package scan

import "encoding/binary"

func Int16(data []byte) (int16, []byte) {
	return int16(binary.BigEndian.Uint16(data)), data[2:]
}

func Uint16(data []byte) (uint16, []byte) {
	return binary.BigEndian.Uint16(data), data[2:]
}

func SetInt16(data []byte, v int16) []byte {
	binary.BigEndian.PutUint16(data, uint16(v))
	return data[2:]
}

func SetUint16(data []byte, v uint16) []byte {
	binary.BigEndian.PutUint16(data, v)
	return data[2:]
}
