package record

import (
	"io"
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
	return WriteAll(w, le.Length, le.Offset)
}

func (le *ListEntry) Read(r io.Reader) error {
	return ReadAll(r, &le.Length, &le.Offset)
}
