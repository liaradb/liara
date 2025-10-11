package record

import (
	"io"

	"github.com/liaradb/liaradb/raw"
)

type ListEntry struct {
	Offset Offset
	Length Offset
}

func newListEntry(offset Offset, length Offset) ListEntry {
	return ListEntry{
		Offset: offset,
		Length: length,
	}
}

func (le ListEntry) Size() int { return le.Offset.Size() + le.Length.Size() }

func (le ListEntry) Write(w io.Writer) error {
	return raw.WriteAll(w, le.Length, le.Offset)
}

func (le *ListEntry) Read(r io.Reader) error {
	return raw.ReadAll(r, &le.Length, &le.Offset)
}
