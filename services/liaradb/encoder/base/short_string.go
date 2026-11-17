package base

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type ShortString struct {
	value string
}

func NewShortString(v string) ShortString {
	return ShortString{v}
}

func (s ShortString) String() string          { return s.value }
func (s ShortString) Bytes() []byte           { return []byte(s.value) }
func (s ShortString) Length() int             { return len(s.value) }
func (s ShortString) Size() int               { return 256 }
func (s ShortString) Write(w io.Writer) error { return raw.WriteShortString(w, s.value) }
func (s *ShortString) Read(r io.Reader) error { return raw.ReadShortString(r, &s.value) }
