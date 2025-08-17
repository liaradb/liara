package raw

import "errors"

var (
	ErrIncompleteWrite = errors.New("incomplete write")
	ErrOverflow        = errors.New("overflow")
	ErrUnderflow       = errors.New("underflow")

	ErrUnableToRead  = errors.New("unable to read")
	ErrUnableToWrite = errors.New("unable to write")
)

func Copy(dst []byte, src []byte) error {
	if !(copy(dst, src) == len(src)) {
		return ErrIncompleteWrite
	}

	return nil
}

// TODO: Is it better to not check offset?
func CopyAt(dst []byte, src []byte, off Offset) error {
	if off < 0 {
		return ErrUnderflow
	}

	if off > Offset(len(dst)) {
		return ErrOverflow
	}

	return Copy(dst[off:], src)
}
