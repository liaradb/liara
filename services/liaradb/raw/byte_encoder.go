package raw

import (
	"encoding/binary"
	"io"
)

type ByteEncoder struct{}

const ByteHeaderSize = 4

func (be *ByteEncoder) Write(w io.Writer, value []byte) error {
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

func (be *ByteEncoder) Read(r io.Reader) ([]byte, error) {
	l, err := be.readLength(r)
	if err != nil {
		return nil, err
	}

	b, err := be.readData(r, l)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (be *ByteEncoder) readLength(r io.Reader) (uint32, error) {
	var l uint32
	err := binary.Read(r, binary.BigEndian, &l)
	return l, err
}

func (be *ByteEncoder) readData(r io.Reader, l uint32) ([]byte, error) {
	b := make([]byte, l)

	if _, err := r.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
