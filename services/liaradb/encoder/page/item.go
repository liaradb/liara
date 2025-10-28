package page

import (
	"io"
	"slices"
)

type Item struct {
	data []byte
}

func NewItem(data []byte) *Item {
	return &Item{data}
}

func NewItemByLength(l Offset) *Item {
	return &Item{make([]byte, l)}
}

func (i *Item) String() string { return string(i.data) }
func (i *Item) Value() []byte  { return i.data } // TODO: Should this clone?
func (i *Item) Length() int    { return len(i.data) }
func (i *Item) Size() int      { return len(i.data) }

func (i *Item) Compare(a *Item) bool {
	return slices.Equal(i.data, a.data)
}

func (i *Item) Write(w io.Writer) error {
	if n, err := w.Write(i.data); err != nil {
		return err
	} else if n < len(i.data) {
		return io.ErrShortWrite
	}

	return nil
}

func (i *Item) Read(r io.Reader) error {
	if _, err := r.Read(i.data); err != nil {
		return err
	}

	return nil
}
