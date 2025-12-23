package mempage

import (
	"io"
)

type FlatList struct {
	length int32
	w      io.OffsetWriter
	r      io.SectionReader
}

type FlatListReader struct {
	r *io.SectionReader
}

func (lr *FlatListReader) Get(i int64) (*FlatListEntry, error) {
	if _, err := lr.r.Seek(FlatListEntrySize*i, io.SeekStart); err != nil {
		return nil, err
	}

	le := &FlatListEntry{}
	return le, le.Read(lr.r)
}

func (lr *FlatListReader) Delete(i int64) error {
	return nil
}
