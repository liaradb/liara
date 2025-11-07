package page

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type ListEntry struct {
	ID     Offset
	Offset Offset
	Length Offset
	CRC    CRC
}

const ListEntrySize = OffsetSize + OffsetSize + ListLengthSize + CrcSize

func newListEntry(id Offset, offset Offset, length Offset) ListEntry {
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
