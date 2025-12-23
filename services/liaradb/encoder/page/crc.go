package page

import (
	"hash/crc32"

	"github.com/liaradb/liaradb/encoder/raw"
)

type CRC struct {
	baseUint32
}

type baseUint32 = raw.BaseUint32

const CrcSize = raw.BaseUint32Size

var table = crc32.MakeTable(crc32.Castagnoli)

func NewCRC(d []byte) CRC {
	return CRC{raw.NewBaseUint32(crc32.Checksum(d, table))}
}

func RestoreCRC[T ~int32 | ~uint32](v T) CRC {
	return CRC{raw.NewBaseUint32(uint32(v))}
}

func (c CRC) Compare(d []byte) bool {
	return NewCRC(d) == c
}
