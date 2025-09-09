package log

import "hash/crc32"

type CRC uint32

var table = crc32.MakeTable(crc32.Castagnoli)

func NewCRC(data []byte) CRC {
	return CRC(crc32.Checksum(data, table))
}

func (c CRC) Compare(data []byte) bool {
	return NewCRC(data) == c
}
