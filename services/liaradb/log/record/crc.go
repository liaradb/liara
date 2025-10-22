package record

import (
	"hash/crc32"
)

type CRC struct {
	baseUint32
}

const CrcSize = baseUint32Size

var table = crc32.MakeTable(crc32.Castagnoli)

func NewCRC(d []byte) CRC {
	return CRC{NewBaseUint32(crc32.Checksum(d, table))}
}

func (c CRC) Compare(d []byte) bool {
	return NewCRC(d) == c
}
