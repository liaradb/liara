package mempage

import (
	"io"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/raw"
)

type listEntry struct {
	ID     page.Offset
	Offset page.Offset
	Length ListLength
	CRC    page.CRC
}

const listEntrySize = page.OffsetSize + page.OffsetSize + ListLengthSize + page.CrcSize

func newListEntry(id page.Offset, offset page.Offset, length ListLength) listEntry {
	return listEntry{
		ID:     id,
		Offset: offset,
		Length: length,
	}
}

func (le listEntry) Size() int { return listEntrySize }

func (le listEntry) Write(w io.Writer) error {
	return raw.WriteAll(w,
		le.ID,
		le.Length,
		le.Offset,
		le.CRC)
}

func (le *listEntry) Read(r io.Reader) error {
	return raw.ReadAll(r,
		&le.ID,
		&le.Length,
		&le.Offset,
		&le.CRC)
}
