package page

import (
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/span"
)

const (
	ItemHeaderSize = page.CrcSize +
		span.FragmentHeaderSize
)

type itemHeader struct {
	crc page.CRC
}
