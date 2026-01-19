package mempage

import (
	"io"
	"slices"

	"github.com/liaradb/liaradb/encoder/page"
)

type item struct {
	data []byte
}

func newItem(data []byte) *item {
	// TODO: Do we need to clone the item?
	return &item{data}
}

func newItemByLength(l ListLength) *item {
	return &item{make([]byte, l)}
}

func (i *item) String() string { return string(i.data) }
func (i *item) Value() []byte  { return i.data } // TODO: Should this clone?
func (i *item) Length() int    { return len(i.data) }
func (i *item) Size() int      { return len(i.data) }

func (i *item) Compare(a *item) bool {
	return slices.Equal(i.data, a.data)
}

func (i *item) Write(w io.Writer) (page.CRC, error) {
	if n, err := w.Write(i.data); err != nil {
		return page.CRC{}, err
	} else if n < len(i.data) {
		return page.CRC{}, io.ErrShortWrite
	}

	return page.NewCRC(i.data), nil
}

func (i *item) Read(r io.Reader, crc page.CRC) error {
	if _, err := r.Read(i.data); err != nil {
		return err
	}

	if !crc.Compare(i.data) {
		return page.ErrInvalidCRC
	}

	return nil
}
