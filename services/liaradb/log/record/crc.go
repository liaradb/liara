package record

import (
	"encoding/binary"
	"hash/crc32"
	"io"
)

type CRC uint32

const CrcSize = 4

var table = crc32.MakeTable(crc32.Castagnoli)

func NewCRC(data []byte) CRC {
	return CRC(crc32.Checksum(data, table))
}

func (c CRC) Compare(data []byte) bool {
	return NewCRC(data) == c
}

func (c CRC) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, c)
}

func (c *CRC) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, c)
}
