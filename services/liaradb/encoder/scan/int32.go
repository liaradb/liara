package scan

import "encoding/binary"

func Int32(data []byte) (int32, []byte, bool) {
	if len(data) < 4 {
		return 0, nil, false
	}

	return int32(binary.BigEndian.Uint32(data[:4])), data[4:], true
}

func Uint32(data []byte) (uint32, []byte, bool) {
	if len(data) < 4 {
		return 0, nil, false
	}

	return binary.BigEndian.Uint32(data[:4]), data[4:], true
}

func SetInt32(data []byte, v int32) ([]byte, bool) {
	if len(data) < 4 {
		return nil, false
	}

	binary.BigEndian.PutUint32(data[:4], uint32(v))
	return data[4:], true
}

func SetUint32(data []byte, v uint32) ([]byte, bool) {
	if len(data) < 4 {
		return nil, false
	}

	binary.BigEndian.PutUint32(data[:4], v)
	return data[4:], true
}
