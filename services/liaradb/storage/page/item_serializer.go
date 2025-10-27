package page

import (
	"io"
	"slices"
)

type ItemSerializer struct {
	data []byte
}

func NewItemSerializer(data []byte) *ItemSerializer {
	return &ItemSerializer{data}
}

func NewItemSerializerByLength(l Offset) *ItemSerializer {
	return &ItemSerializer{make([]byte, l)}
}

func (is *ItemSerializer) String() string { return string(is.data) }
func (is *ItemSerializer) Value() []byte  { return is.data } // TODO: Should this clone?
func (is *ItemSerializer) Length() int    { return len(is.data) }
func (is *ItemSerializer) Size() int      { return len(is.data) }

func (is *ItemSerializer) Compare(a *ItemSerializer) bool {
	return slices.Equal(is.data, a.data)
}

func (is *ItemSerializer) Write(w io.Writer) error {
	if n, err := w.Write(is.data); err != nil {
		return err
	} else if n < len(is.data) {
		return io.ErrShortWrite
	}

	return nil
}

func (is *ItemSerializer) Read(r io.Reader) error {
	if _, err := r.Read(is.data); err != nil {
		return err
	}

	return nil
}
