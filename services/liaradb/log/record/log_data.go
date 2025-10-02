package record

import (
	"encoding/binary"
	"io"
)

type LogData struct {
	data []byte
}

const LogDataHeaderSize = 4

func NewLogData(data []byte) *LogData {
	return &LogData{
		data: data,
	}
}

func (ld *LogData) Bytes() []byte { return ld.data }
func (ld *LogData) Length() int   { return len(ld.data) }

func (ld LogData) Size() int {
	return ld.Length() + LogDataHeaderSize
}

func (ld *LogData) Write(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, uint32(len(ld.data))); err != nil {
		return err
	}

	_, err := w.Write(ld.data)
	return err
}

func (ld *LogData) Read(r io.Reader) error {
	l, err := ld.readLength(r)
	if err != nil {
		return err
	}

	b, err := ld.readData(r, l)
	if err != nil {
		return err
	}

	ld.data = b
	return nil
}

func (ld *LogData) readLength(r io.Reader) (uint32, error) {
	var l uint32
	err := binary.Read(r, binary.BigEndian, &l)
	return l, err
}

func (ld *LogData) readData(r io.Reader, l uint32) ([]byte, error) {
	b := make([]byte, l)

	if _, err := r.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
