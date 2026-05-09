package page

import "github.com/liaradb/liaradb/encoder/page"

const (
	itemHeaderSize = page.CrcSize
)

type itemHeader struct {
	crc page.CRC
}
