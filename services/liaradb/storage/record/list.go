package record

import (
	"io"
)

type List struct {
	entries []ListEntry
}

func (l List) Length() ListLength { return ListLength(len(l.entries)) }

func (l *List) Add(offset Offset, length Offset) (int, error) {
	// TODO: Test this
	if int(offset) < l.space() {
		return 0, ErrInsufficientSpace
	}

	le := newListEntry(offset, length)
	l.entries = append(l.entries, le)
	return len(l.entries) - 1, nil
}

func (l List) space() int {
	return l.Size() + ListEntrySize
}

func (l List) Size() int {
	s := ListLength(0).Size()
	for _, e := range l.entries {
		s += e.Size()
	}
	return s
}

func (l List) offset(index int) Offset {
	if index < 0 || index >= len(l.entries) {
		return 0
	}

	return l.entries[index].Offset
}

func (l List) entriesSize() int {
	var s int
	for _, e := range l.entries {
		s += e.Size()
	}
	return s
}

func (l List) Write(w io.Writer) error {
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

func (l *List) Read(r io.Reader) error {
	var length ListLength
	if err := length.Read(r); err != nil {
		return err
	}

	entries := make([]ListEntry, 0, length)
	for range length {
		var le ListEntry
		if err := le.Read(r); err != nil {
			return err
		}

		entries = append(entries, le)
	}

	l.entries = entries
	return nil
}
