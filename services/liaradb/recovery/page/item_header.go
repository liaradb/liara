package page

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

const (
	ItemHeaderSize = page.CrcSize +
		record.FragmentHeaderSize
)

type itemHeader struct {
	crc page.CRC
}
