package raw

import "encoding/binary"

// TODO: Is it better to just use []byte, and not string?

func CopyString(dst []byte, value string, off Offset) error {
	if off < 0 {
		return ErrUnderflow
	}

	strOff := off + stringHeaderOffset
	strLen := Offset(len(value))

	if strOff+strLen > Offset(len(dst)) {
		return ErrOverflow
	}

	binary.BigEndian.PutUint32(dst[off:], uint32(strLen))

	if !(copy(dst[strOff:], []byte(value)) == int(strLen)) {
		return ErrIncompleteWrite
	}

	return nil
}

func GetString(src []byte, off Offset) (string, error) {
	// TODO: Should there be a maximum string length?
	strLen, err := GetUint32(src, off)
	if err != nil {
		return "", err
	}

	strOff := off + stringHeaderOffset
	l := Offset(strLen)
	if strOff+l > Offset(len(src)) {
		return "", ErrOverflow
	}

	return string(src[strOff : strOff+l]), nil
}
