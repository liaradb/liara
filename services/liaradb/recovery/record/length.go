package record

import "github.com/liaradb/liaradb/encoder/base"

type Length struct {
	baseUint32
}

const LengthSize = 4

func NewLength(size uint32) Length {
	return Length{base.NewBaseUint32(size)}
}
