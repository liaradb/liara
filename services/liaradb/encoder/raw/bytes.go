package raw

import (
	"io"
	"slices"
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
func (b *Bytes) Size() int               { return len(b.data) + HeaderSize }
func (b *Bytes) Write(w io.Writer) error { return Write(w, b.data) }
func (b *Bytes) Read(r io.Reader) error  { return Read(r, &b.data) }
