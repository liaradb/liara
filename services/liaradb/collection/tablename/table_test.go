package tablename

import (
	"fmt"
	"testing"

	"github.com/liaradb/liaradb/domain/value"
)

func TestTable(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	tf := NewTenantFolder(tid)

	name := "default"
	tbl := tf.Table(name)

	col := tbl.Collection()

	colWant := fmt.Sprintf("%v/%v.col", tid, name)
	if s := col.String(); s != colWant {
		t.Errorf("incorrect string: %v, expected: %v", s, colWant)
	}

	idx := tbl.Index(1, value.NewPartitionID(2))

	idxWant := fmt.Sprintf("%v/%v--%08x--%08x.bt", tid, name, 1, 2)
	if s := idx.String(); s != idxWant {
		t.Errorf("incorrect string: %v, expected: %v", s, idxWant)
	}
}
