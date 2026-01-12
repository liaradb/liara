package scan

import "encoding/binary"

func Int32(data []byte) (int32, []byte) {
	return int32(binary.BigEndian.Uint32(data[:4])), data[4:]
}

func Uint32(data []byte) (uint32, []byte) {
	return binary.BigEndian.Uint32(data[:4]), data[4:]
}

func SetInt32(data []byte, v int32) []byte {
	binary.BigEndian.PutUint32(data[:4], uint32(v))
	return data[4:]
}

func SetUint32(data []byte, v uint32) []byte {
	binary.BigEndian.PutUint32(data[:4], v)
	return data[4:]
}
