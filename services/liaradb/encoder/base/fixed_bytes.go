package base

import (
	"io"
	"slices"
)

type FixedBytes struct {
	data []byte
}

func NewFixedBytes(data []byte) *FixedBytes {
	return &FixedBytes{data}
}

func NewFixedBytesFromSize(size int) *FixedBytes {
	return &FixedBytes{make([]byte, size)}
}

func (b *FixedBytes) String() string             { return string(b.data) }
func (b *FixedBytes) Compare(a *FixedBytes) bool { return slices.Equal(b.data, a.data) }
func (b *FixedBytes) Value() []byte              { return b.data }
func (b *FixedBytes) Length() int                { return len(b.data) }
func (b *FixedBytes) Size() int                  { return len(b.data) }

func (b *FixedBytes) Write(w io.Writer) error {
	_, err := w.Write(b.data)
	return err
}

func (b *FixedBytes) Read(r io.Reader) error {
	_, err := r.Read(b.data)
	return err
}
