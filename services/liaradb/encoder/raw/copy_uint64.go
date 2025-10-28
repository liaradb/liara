package raw

import "encoding/binary"

func CopyUint64(dst []byte, value uint64, off Offset) error {
	if off < 0 {
		return ErrUnderflow
	}

	if off+Uint64Length > Offset(len(dst)) {
		return ErrOverflow
	}

	binary.BigEndian.PutUint64(dst[off:], value)

	return nil
}

func GetUint64(src []byte, off Offset) (uint64, error) {
	if off < 0 {
		return 0, ErrUnderflow
	}

	if off+Uint64Length > Offset(len(src)) {
		return 0, ErrOverflow
	}

	return binary.BigEndian.Uint64(src[off:]), nil
}

func CopyInt64(dst []byte, value, int64, off Offset) error {
	return CopyUint64(dst, uint64(value), off)
}

func GetInt64(src []byte, off Offset) (int64, error) {
	v, err := GetUint64(src, off)
	return int64(v), err
}
