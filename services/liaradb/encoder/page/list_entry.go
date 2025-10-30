package page

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

// TODO: Include a Serializer?

type ListEntry struct {
	Offset Offset
	Length Offset
	CRC    CRC
}

const ListEntrySize = OffsetSize + ListLengthSize + CrcSize

func newListEntry(offset Offset, length Offset) ListEntry {
	return ListEntry{
		Offset: offset,
		Length: length,
	}
}

func (le ListEntry) Size() int { return ListEntrySize }

func (le ListEntry) Write(w io.Writer) error {
	return raw.WriteAll(w, le.Length, le.Offset, le.CRC)
}

func (le *ListEntry) Read(r io.Reader) error {
	return raw.ReadAll(r, &le.Length, &le.Offset, &le.CRC)
}
