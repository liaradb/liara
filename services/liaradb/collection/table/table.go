package table

import "github.com/liaradb/liaradb/collection/table/schema"

type Table struct {
	id       ID
	schemaID schema.ID
}

func New(
	id ID,
	schemaID schema.ID,
) *Table {
	return &Table{
		id:       id,
		schemaID: schemaID,
	}
}

func (t *Table) ID() ID              { return t.id }
func (t *Table) SchemaID() schema.ID { return t.schemaID }
