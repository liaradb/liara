package page

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type FlatListEntry struct {
	ID     Offset
	Offset Offset
	Length ListLength
	CRC    CRC
}

const FlatListEntrySize = OffsetSize + OffsetSize + ListLengthSize + CrcSize

func newFlatListEntry(id Offset, offset Offset, length ListLength) ListEntry {
	return ListEntry{
		ID:     id,
		Offset: offset,
		Length: length,
	}
}

func (le FlatListEntry) Size() int { return ListEntrySize }

func (le FlatListEntry) Write(w io.Writer) error {
	return raw.WriteAll(w,
		le.ID,
		le.Length,
		le.Offset,
		le.CRC)
}

func (le *FlatListEntry) Read(r io.Reader) error {
	return raw.ReadAll(r,
		&le.ID,
		&le.Length,
		&le.Offset,
		&le.CRC)
}
