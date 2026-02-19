package value

import (
	"io"
	"slices"

	"github.com/liaradb/liaradb/encoder/raw"
)

type Data struct {
	data []byte
}

func NewData(data []byte) Data {
	return Data{
		data: data,
	}
}

func (d Data) String() string { return string(d.data) }
func (d Data) Value() []byte  { return d.data }
func (d Data) Size() int      { return len(d.data) }

func (d Data) Write(w io.Writer) error {
	return raw.Write(w, d.data)
}

func (d *Data) Read(r io.Reader) error {
	return raw.Read(r, &d.data)
}

func (d *Data) Compare(b *Data) bool {
	return slices.Equal(d.data, b.data)
}
