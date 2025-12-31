package mempage

import (
	"io"
	"slices"

	"github.com/liaradb/liaradb/encoder/page"
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

func (i *Item) Write(w io.Writer) (page.CRC, error) {
	if n, err := w.Write(i.data); err != nil {
		return page.CRC{}, err
	} else if n < len(i.data) {
		return page.CRC{}, io.ErrShortWrite
	}

	return page.NewCRC(i.data), nil
}

func (i *Item) Read(r io.Reader, crc page.CRC) error {
	if _, err := r.Read(i.data); err != nil {
		return err
	}

	if !crc.Compare(i.data) {
		return page.ErrInvalidCRC
	}

	return nil
}
