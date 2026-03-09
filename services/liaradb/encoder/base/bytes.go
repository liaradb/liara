package base

import (
	"io"
	"slices"

	"github.com/liaradb/liaradb/encoder/raw"
)

type Bytes struct {
	data []byte
}

func NewBytes(data []byte) *Bytes {
	return &Bytes{data}
}

func (b *Bytes) String() string          { return string(b.data) }
func (b *Bytes) Compare(a *Bytes) bool   { return slices.Equal(b.data, a.data) }
func (b *Bytes) Value() []byte           { return b.data } // TODO: Should this clone?
func (b *Bytes) Length() int             { return len(b.data) }
func (b *Bytes) Size() int               { return len(b.data) + raw.HeaderSize }
func (b *Bytes) Write(w io.Writer) error { return raw.Write(w, b.data) }
func (b *Bytes) Read(r io.Reader) error  { return raw.Read(r, &b.data) }

func (b Bytes) WriteData(data []byte, colSize int) []byte {
	data[0] = byte(len(b.data))
	copy(data[1:colSize], []byte(b.data))
	clear(data[1+len(b.data) : colSize])
	return data[colSize:]
}

func (b *Bytes) ReadData(data []byte, colSize int) []byte {
	l := data[0]
	b.data = data[1 : 1+l]
	return data[colSize:]
}
