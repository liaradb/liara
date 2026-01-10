package table

import "github.com/liaradb/liaradb/collection/schema"

type Table struct {
	id       ID
	schemaID schema.ID
}

func (t *Table) ID() ID              { return t.id }
func (t *Table) SchemaID() schema.ID { return t.schemaID }
