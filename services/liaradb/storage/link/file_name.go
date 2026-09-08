package link

import "github.com/liaradb/liaradb/encoder/raw"

type FileName string

func NewFileName(value string) FileName {
	return FileName(value)
}

func (fn FileName) String() string { return string(fn) }

// TODO: Verify this size
func (fn FileName) Size() int { return raw.StringSize(fn) }

func (fn FileName) BlockID(position FilePosition) BlockID {
	return NewBlockID(fn, position)
}
