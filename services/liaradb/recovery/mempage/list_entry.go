package mempage

import (
	"io"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/raw"
)

type ListEntry struct {
	ID     page.Offset
	Offset page.Offset
	Length ListLength
	CRC    page.CRC
}

const ListEntrySize = page.OffsetSize + page.OffsetSize + ListLengthSize + page.CrcSize

func newListEntry(id page.Offset, offset page.Offset, length ListLength) ListEntry {
	return ListEntry{
		ID:     id,
		Offset: offset,
		Length: length,
	}
}

func (le ListEntry) Size() int { return ListEntrySize }

func (le ListEntry) Write(w io.Writer) error {
	return raw.WriteAll(w,
		le.ID,
		le.Length,
		le.Offset,
		le.CRC)
}

func (le *ListEntry) Read(r io.Reader) error {
	return raw.ReadAll(r,
		&le.ID,
		&le.Length,
		&le.Offset,
		&le.CRC)
}
