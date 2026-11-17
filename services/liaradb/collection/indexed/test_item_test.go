package indexed

import (
	"io"

	"github.com/liaradb/liaradb/encoder/serializer"
)

type TestItem struct {
	id   TestID
	name TestName
}

func NewTestItem(
	id TestID,
	name TestName,
) *TestItem {
	return &TestItem{
		id:   id,
		name: name,
	}
}

func RestoreTestItem(
	id TestID,
	name TestName,
) *TestItem {
	return &TestItem{
		id:   id,
		name: name,
	}
}

func (ti *TestItem) ID() TestID     { return ti.id }
func (ti *TestItem) Name() TestName { return ti.name }

func (ti *TestItem) Compare(b *TestItem) bool {
	if ti == b {
		return true
	}

	return ti.id == b.id
}

func (ti *TestItem) Size() int {
	return serializer.Size(
		ti.id,
		ti.name)
}

func (ti *TestItem) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		ti.id,
		ti.name)
}

func (ti *TestItem) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&ti.id,
		&ti.name)
}
