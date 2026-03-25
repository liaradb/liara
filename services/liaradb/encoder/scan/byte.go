package scan

func Byte(data []byte) (byte, []byte, bool) {
	if len(data) < 1 {
		return 0, nil, false
	}

	return data[0], data[1:], true
}

func Int8(data []byte) (int8, []byte, bool) {
	if len(data) < 1 {
		return 0, nil, false
	}

	return int8(data[0]), data[1:], true
}

func SetByte(data []byte, v byte) ([]byte, bool) {
	if len(data) < 1 {
		return nil, false
	}

	data[0] = v
	return data[1:], true
}

func SetInt8(data []byte, v int8) ([]byte, bool) {
	if len(data) < 1 {
		return nil, false
	}

	data[0] = byte(v)
	return data[1:], true
}
