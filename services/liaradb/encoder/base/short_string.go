package base

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type ShortString[S ShortStringSizer] struct {
	value string
}

type ShortStringSizer interface {
	StringSize() byte
}

func NewShortString[S ShortStringSizer](v string) ShortString[S] {
	return ShortString[S]{v}
}

func (s ShortString[S]) String() string          { return s.value }
func (s ShortString[S]) Bytes() []byte           { return []byte(s.value) }
func (s ShortString[S]) Length() int             { return len(s.value) }
func (s ShortString[S]) Size() int               { var ss S; return int(ss.StringSize()) + 1 }
func (s ShortString[S]) Write(w io.Writer) error { return raw.WriteShortString(w, s.value) }
func (s *ShortString[S]) Read(r io.Reader) error { return raw.ReadShortString(r, &s.value) }

type ShortString255 struct {
}

func (ss ShortString255) StringSize() byte { return 255 }
