package page

import (
	"hash/crc32"

	"github.com/liaradb/liaradb/encoder/base"
)

type CRC struct {
	baseUint32
}

type baseUint32 = base.Uint32

const CrcSize = base.Uint32Size

var table = crc32.MakeTable(crc32.Castagnoli)

func NewCRC(d []byte) CRC {
	return CRC{base.NewUint32(crc32.Checksum(d, table))}
}

func RestoreCRC[T ~int32 | ~uint32](v T) CRC {
	return CRC{base.NewUint32(uint32(v))}
}

func (c CRC) Compare(d []byte) bool {
	return NewCRC(d) == c
}
