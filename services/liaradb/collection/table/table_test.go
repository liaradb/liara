package table

import (
	"testing"

	"github.com/liaradb/liaradb/collection/table/schema"
)

func TestTable(t *testing.T) {
	id := NewID()
	sid := schema.NewID()

	tbl := New(id, sid)

	if v := tbl.ID(); v != id {
		t.Errorf("incorrect id: %v, expected: %v", v, id)
	}

	if v := tbl.SchemaID(); v != sid {
		t.Errorf("incorrect schema: %v, expected: %v", v, sid)
	}
}
