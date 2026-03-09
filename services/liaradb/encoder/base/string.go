package base

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type String string

func (s String) String() string          { return string(s) }
func (s String) Bytes() []byte           { return []byte(s) }
func (s String) Length() int             { return len(s) }
func (s String) Size() int               { return raw.StringSize(s) }
func (s String) Write(w io.Writer) error { return raw.WriteString(w, s) }
func (s *String) Read(r io.Reader) error { return raw.ReadString(r, s) }

func (s String) WriteData(data []byte, colSize int) []byte {
	data[0] = byte(len(s))
	copy(data[1:colSize], []byte(s))
	clear(data[1+len(s) : colSize])
	return data[colSize:]
}

func (s *String) ReadData(data []byte, colSize int) []byte {
	l := data[0]
	*s = String(data[1 : 1+l])
	return data[colSize:]
}
