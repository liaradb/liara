package entity

import (
	"io"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
)

type Row struct {
	ID          value.RowID
	Version     value.Version
	PartitionID value.PartitionID
	Name        value.RowName
	Schema      value.Schema
	Metadata    Metadata
	Data        value.Data
}

// TODO: Test this
func (e Row) Size() int {
	return raw.Size(
		e.ID,
		e.Version,
		e.PartitionID,
		e.Name,
		e.Schema,
		e.Metadata,
		e.Data)
}

func (e Row) Write(w io.Writer) error {
	return raw.WriteAll(w,
		e.ID,
		e.Version,
		e.PartitionID,
		e.Name,
		e.Schema,
		e.Metadata,
		e.Data)
}

func (e *Row) Read(r io.Reader) error {
	return raw.ReadAll(r,
		&e.ID,
		&e.Version,
		&e.PartitionID,
		&e.Name,
		&e.Schema,
		&e.Metadata,
		&e.Data)
}
