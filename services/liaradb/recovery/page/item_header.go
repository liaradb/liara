package page

import "github.com/liaradb/liaradb/encoder/page"

const (
	ItemHeaderSize = page.CrcSize
)

type itemHeader struct {
	crc page.CRC
}
