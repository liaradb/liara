package raw

import "encoding/binary"

func CopyUint32(dst []byte, value uint32, off Offset) error {
	if off < 0 {
		return ErrUnderflow
	}

	if off+Uint32Length > Offset(len(dst)) {
		return ErrOverflow
	}

	binary.BigEndian.PutUint32(dst[off:], value)

	return nil
}

func GetUint32(src []byte, off Offset) (uint32, error) {
	if off < 0 {
		return 0, ErrUnderflow
	}

	if off+Uint32Length > Offset(len(src)) {
		return 0, ErrOverflow
	}

	return binary.BigEndian.Uint32(src[off:]), nil
}

func CopyInt32(dst []byte, value, int32, off Offset) error {
	return CopyUint32(dst, uint32(value), off)
}

func GetInt32(src []byte, off Offset) (int32, error) {
	v, err := GetUint32(src, off)
	return int32(v), err
}
