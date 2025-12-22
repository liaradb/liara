package link

import "github.com/liaradb/liaradb/encoder/page"

type FileName string

func NewFileName(value string) FileName {
	return FileName(value)
}

func (fn FileName) String() string { return string(fn) }

func (fn FileName) BlockID(position page.Offset) BlockID {
	return NewBlockID(fn.String(), position)
}
