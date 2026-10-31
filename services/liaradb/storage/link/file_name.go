package link

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type FileName string

func NewFileName(value string) FileName {
	return FileName(value)
}

func (fn FileName) String() string { return string(fn) }

func (fn FileName) Size() int               { return raw.StringSize(fn) }
func (fn FileName) Write(w io.Writer) error { return raw.WriteString(w, fn) }
func (fn *FileName) Read(r io.Reader) error { return raw.ReadString(r, fn) }

func (fn FileName) FileBlockID(bid BlockID) FileBlockID {
	return NewFileBlockID(fn, bid)
}
