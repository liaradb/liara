package record

import (
	"io"
)

type List struct {
	entries []ListEntry
}

func (l List) Length() ListLength { return ListLength(len(l.entries)) }

func (l *List) Add(length Offset) int {
	le := newListEntry(Offset(l.Size()), length)
	l.entries = append(l.entries, le)
	return len(l.entries) - 1
}

func (l List) Size() int {
	s := ListLength(0).Size()
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
