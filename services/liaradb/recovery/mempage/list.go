package mempage

import (
	"io"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/raw"
)

type list struct {
	headerSize int
	highWater  page.Offset
	entries    []listEntry
}

// TODO: Test headersize
func newList(headerSize int) list {
	return list{
		headerSize: headerSize,
	}
}

func (l *list) Reset() {
	l.highWater = 0
	l.entries = nil
}

func (l list) Length() listLength { return listLength(len(l.entries)) }

func (l *list) Add(offset page.Offset, length listLength) (page.Offset, error) {
	// TODO: Test this
	if int(offset) < l.space() {
		return 0, raw.ErrInsufficientSpace
	}

	id := l.highWater
	le := newListEntry(id, offset, length)
	l.highWater++
	l.entries = append(l.entries, le)
	return id, nil
}

func (l list) space() int {
	return l.Size() + listEntrySize
}

func (l list) Size() int {
	s := page.Offset(0).Size() + listLength(0).Size() + l.headerSize
	for _, e := range l.entries {
		s += e.Size()
	}
	return s
}

func (l list) offset(index int) page.Offset {
	if index < 0 || index >= len(l.entries) {
		return 0
	}

	return l.entries[index].Offset
}

func (l *list) setCRC(index int, crc page.CRC) {
	if index < 0 || index >= len(l.entries) {
		return
	}

	l.entries[index].CRC = crc
}

func (l list) entriesSize() int {
	var s int
	for _, e := range l.entries {
		s += e.Length.Value()
	}
	return s
}

func (l list) Write(w io.Writer) error {
	if err := l.highWater.Write(w); err != nil {
		return err
	}

	if err := l.Length().Write(w); err != nil {
		return err
	}

	for _, e := range l.entries {
		if err := e.Write(w); err != nil {
			return err
		}
	}

	return nil
}

func (l *list) Read(r io.Reader) error {
	if err := l.highWater.Read(r); err != nil {
		return err
	}

	var length listLength
	if err := length.Read(r); err != nil {
		return err
	}

	entries := make([]listEntry, 0, length)
	for range length {
		var le listEntry
		if err := le.Read(r); err != nil {
			return err
		}

		entries = append(entries, le)
	}

	l.entries = entries
	return nil
}
