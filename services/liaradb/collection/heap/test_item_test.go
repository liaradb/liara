package heap

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

const (
	RequestLogSize = 0 +
		TestIDSize
)

type TestItem struct {
	id TestID
}

func NewTestItem(
	id TestID,
) *TestItem {
	return &TestItem{
		id: id,
	}
}

func RestoreTestItem(
	id TestID,
) *TestItem {
	return &TestItem{
		id: id,
	}
}

func (ti *TestItem) ID() TestID { return ti.id }
func (ti *TestItem) Size() int  { return RequestLogSize }

func (ti *TestItem) Compare(b *TestItem) bool {
	if ti == b {
		return true
	}

	return ti.id == b.id
}

func (ti *TestItem) WriteData(data []byte) []byte {
	return ti.id.WriteData(data)
}

func (ti *TestItem) ReadData(data []byte) []byte {
	return ti.id.ReadData(data)
}

func (ti *TestItem) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		ti.id)
}

func (ti *TestItem) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&ti.id)
}
