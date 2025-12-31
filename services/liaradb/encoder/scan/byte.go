package scan

func Byte(data []byte) (byte, []byte) {
	return data[0], data[1:]
}

func Int8(data []byte) (int8, []byte) {
	return int8(data[0]), data[1:]
}

func SetByte(data []byte, v byte) []byte {
	data[0] = v
	return data[1:]
}

func SetInt8(data []byte, v int8) []byte {
	data[0] = byte(v)
	return data[1:]
}
