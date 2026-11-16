package heap

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

type TestItem struct {
	name TestName
}

func NewTestItem(
	id TestName,
) *TestItem {
	return &TestItem{
		name: id,
	}
}

func RestoreTestItem(
	name TestName,
) *TestItem {
	return &TestItem{
		name: name,
	}
}

func (ti *TestItem) Name() TestName { return ti.name }
func (ti *TestItem) Size() int      { return serializer.Size(ti.name) }

func (ti *TestItem) Compare(b *TestItem) bool {
	if ti == b {
		return true
	}

	return ti.name == b.name
}

func (ti *TestItem) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		ti.name)
}

func (ti *TestItem) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&ti.name)
}
