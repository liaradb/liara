package page

import (
	"hash/crc32"

	"github.com/liaradb/liaradb/encoder/base"
)

type CRC struct {
	baseUint32
}

type baseUint32 = base.BaseUint32

const CrcSize = base.BaseUint32Size

var table = crc32.MakeTable(crc32.Castagnoli)

func NewCRC(d []byte) CRC {
	return CRC{base.NewBaseUint32(crc32.Checksum(d, table))}
}

func RestoreCRC[T ~int32 | ~uint32](v T) CRC {
	return CRC{base.NewBaseUint32(uint32(v))}
}

func (c CRC) Compare(d []byte) bool {
	return NewCRC(d) == c
}
