package scan

import "encoding/binary"

func Int64(data []byte) (int64, []byte, bool) {
	if len(data) < 8 {
		return 0, nil, false
	}

	return int64(binary.BigEndian.Uint64(data[:8])), data[8:], true
}

func Uint64(data []byte) (uint64, []byte, bool) {
	if len(data) < 8 {
		return 0, nil, false
	}

	return binary.BigEndian.Uint64(data[:8]), data[8:], true
}

func SetInt64(data []byte, v int64) ([]byte, bool) {
	if len(data) < 8 {
		return nil, false
	}

	binary.BigEndian.PutUint64(data[:8], uint64(v))
	return data[8:], true
}

func SetUint64(data []byte, v uint64) ([]byte, bool) {
	if len(data) < 8 {
		return nil, false
	}

	binary.BigEndian.PutUint64(data[:8], v)
	return data[8:], true
}
