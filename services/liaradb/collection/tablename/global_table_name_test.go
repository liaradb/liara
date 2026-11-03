package tablename

import (
	"fmt"
	"testing"

	"github.com/liaradb/liaradb/domain/value"
)

func TestGlobalTable(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	tf := NewTenantFolder(tid)

	name := "default"
	tbl := tf.GlobalTable(name)

	col := tbl.Collection()

	colWant := fmt.Sprintf("%v/%v.col", tid, name)
	if s := col.String(); s != colWant {
		t.Errorf("incorrect string: %v, expected: %v", s, colWant)
	}

	idx := tbl.Index(1)

	idxWant := fmt.Sprintf("%v/%v--%08x.bt", tid, name, 1)
	if s := idx.String(); s != idxWant {
		t.Errorf("incorrect string: %v, expected: %v", s, idxWant)
	}
}
