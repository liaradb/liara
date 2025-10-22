package record

type Length struct {
	baseUint32
}

const LengthSize = 4

func NewLength(size uint32) Length {
	return Length{NewBaseUint32(size)}
}
