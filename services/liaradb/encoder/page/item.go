package page

import (
	"io"
	"slices"
)

type Item struct {
	data []byte
}

func NewItem(data []byte) *Item {
	// TODO: Do we need to clone the item?
	return &Item{data}
}

func NewItemByLength(l ListLength) *Item {
	return &Item{make([]byte, l)}
}

func (i *Item) String() string { return string(i.data) }
func (i *Item) Value() []byte  { return i.data } // TODO: Should this clone?
func (i *Item) Length() int    { return len(i.data) }
func (i *Item) Size() int      { return len(i.data) }

func (i *Item) Compare(a *Item) bool {
	return slices.Equal(i.data, a.data)
}

func (i *Item) Write(w io.Writer) (CRC, error) {
	if n, err := w.Write(i.data); err != nil {
		return CRC{}, err
	} else if n < len(i.data) {
		return CRC{}, io.ErrShortWrite
	}

	return NewCRC(i.data), nil
}

func (i *Item) Read(r io.Reader, crc CRC) error {
	if _, err := r.Read(i.data); err != nil {
		return err
	}

	if !crc.Compare(i.data) {
		return ErrInvalidCRC
	}

	return nil
}
