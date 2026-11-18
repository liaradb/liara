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
