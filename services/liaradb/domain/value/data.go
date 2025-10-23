package value

import (
	"io"

	"github.com/liaradb/liaradb/raw"
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

func (d Data) Write(w io.Writer) error {
	return raw.Write(w, d.data)
}

func (d *Data) Read(r io.Reader) error {
	return raw.Read(r, &d.data)
}
