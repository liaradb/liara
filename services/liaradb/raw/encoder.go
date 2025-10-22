package raw

import (
	"encoding/binary"
	"io"
)

const HeaderSize = 4

// TODO: Can we rename this Size?
func StringSize[S ~string](s S) int {
	return HeaderSize + len(s)
}

func Write(w io.Writer, value []byte) error {
	if err := binary.Write(w, binary.BigEndian, uint32(len(value))); err != nil {
		return err
	}

	if n, err := w.Write([]byte(value)); err != nil {
		return err
	} else if n < len(value) {
		return io.ErrShortWrite
	}

	return nil
}

// TODO: Can we reuse value?
func Read(r io.Reader, value *[]byte) error {
	l, err := readLength(r)
	if err != nil {
		return err
	}

	b, err := readData(r, l)
	if err != nil {
		return err
	}

	*value = b
	return nil
}

func WriteString[S ~string](w io.Writer, value S) error {
	if err := binary.Write(w, binary.BigEndian, uint32(len(value))); err != nil {
		return err
	}

	if n, err := w.Write([]byte(value)); err != nil {
		return err
	} else if n < len(value) {
		return io.ErrShortWrite
	}

	return nil
}

func ReadString[S ~string](r io.Reader, s *S) error {
	l, err := readLength(r)
	if err != nil {
		return err
	}

	b, err := readData(r, l)
	if err != nil {
		return err
	}

	*s = S(b)
	return nil
}

func readLength(r io.Reader) (uint32, error) {
	var l uint32
	err := binary.Read(r, binary.BigEndian, &l)
	return l, err
}

func readData(r io.Reader, l uint32) ([]byte, error) {
	b := make([]byte, l)

	if _, err := r.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
