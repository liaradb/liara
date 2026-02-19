package entity

import (
	"io"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
)

type Row struct {
	id          value.RowID
	version     value.Version
	partitionID value.PartitionID
	name        value.RowName
	schema      value.Schema
	metadata    Metadata
	data        value.Data
}

func NewRow(
	id value.RowID,
	version value.Version,
	partitionID value.PartitionID,
	name value.RowName,
	schema value.Schema,
	metadata Metadata,
	data value.Data,
) *Row {
	return &Row{
		id:          id,
		version:     version,
		partitionID: partitionID,
		name:        name,
		schema:      schema,
		metadata:    metadata,
		data:        data,
	}
}

func (r *Row) ID() value.RowID                { return r.id }
func (r *Row) Version() value.Version         { return r.version }
func (r *Row) PartitionID() value.PartitionID { return r.partitionID }
func (r *Row) Name() value.RowName            { return r.name }
func (r *Row) Schema() value.Schema           { return r.schema }
func (r *Row) Metadata() Metadata             { return r.metadata }
func (r *Row) Data() value.Data               { return r.data }

func (r *Row) SetData(data value.Data) {
	r.data = data
}

// TODO: Test this
func (r *Row) Size() int {
	return raw.Size(
		r.id,
		r.version,
		r.partitionID,
		r.name,
		r.schema,
		r.metadata,
		r.data)
}

func (r *Row) Write(w io.Writer) error {
	return raw.WriteAll(w,
		r.id,
		r.version,
		r.partitionID,
		r.name,
		r.schema,
		r.metadata,
		r.data)
}

func (e *Row) Read(r io.Reader) error {
	return raw.ReadAll(r,
		&e.id,
		&e.version,
		&e.partitionID,
		&e.name,
		&e.schema,
		&e.metadata,
		&e.data)
}

func (e *Row) Compare(b *Row) bool {
	if e == b {
		return true
	}

	if e.id != b.id ||
		e.metadata != b.metadata ||
		e.name != b.name ||
		e.partitionID != b.partitionID ||
		e.schema != b.schema ||
		e.version != b.version {
		return false
	}

	return e.data.Compare(&b.data)
}
