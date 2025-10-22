package record

import "github.com/liaradb/liaradb/raw"

type Length struct {
	baseUint32
}

const LengthSize = 4

func NewLength(size uint32) Length {
	return Length{raw.NewBaseUint32(size)}
}
