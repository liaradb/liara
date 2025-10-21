package raw

import (
	"encoding/binary"
	"io"
)

type StringEncoder struct{}

const StringHeaderSize = 4

func (se *StringEncoder) Write(w io.Writer, value string) error {
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

func (se *StringEncoder) Read(r io.Reader) (string, error) {
	l, err := se.readLength(r)
	if err != nil {
		return "", err
	}

	b, err := se.readData(r, l)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (se *StringEncoder) readLength(r io.Reader) (uint32, error) {
	var l uint32
	err := binary.Read(r, binary.BigEndian, &l)
	return l, err
}

func (se *StringEncoder) readData(r io.Reader, l uint32) ([]byte, error) {
	b := make([]byte, l)

	if _, err := r.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
