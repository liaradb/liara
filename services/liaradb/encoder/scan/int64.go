package scan

import "encoding/binary"

func Int64(data []byte) (int64, []byte) {
	return int64(binary.BigEndian.Uint64(data)), data[8:]
}

func Uint64(data []byte) (uint64, []byte) {
	return binary.BigEndian.Uint64(data), data[8:]
}

func SetInt64(data []byte, v int64) []byte {
	binary.BigEndian.PutUint64(data, uint64(v))
	return data[8:]
}

func SetUint64(data []byte, v uint64) []byte {
	binary.BigEndian.PutUint64(data, v)
	return data[8:]
}
