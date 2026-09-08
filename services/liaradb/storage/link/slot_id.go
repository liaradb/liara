package link

import (
	"fmt"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/scan"
)

type SlotID int16

const SlotIDSize = 2

func (i SlotID) Value() int16   { return int16(i) }
func (i SlotID) Size() int      { return SlotIDSize }
func (i SlotID) String() string { return fmt.Sprintf("%v", i.Value()) }

func (i SlotID) Write(w io.Writer) error {
	return raw.WriteInt16(w, i)
}

func (i *SlotID) Read(r io.Reader) error {
	return raw.ReadInt16(r, i)
}

func (i SlotID) WriteData(data []byte) ([]byte, bool) {
	return scan.SetInt16(data, i.Value())
}

func (i *SlotID) ReadData(data []byte) ([]byte, bool) {
	block, data0, ok := scan.Int16(data)
	if !ok {
		return nil, false
	}

	*i = SlotID(block)
	return data0, true
}
