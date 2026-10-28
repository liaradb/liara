package transaction

import (
	"fmt"
	"testing"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/value"
)

func TestItemID(t *testing.T) {
	t.Parallel()

	id := "id"
	tid := value.NewTenantID()
	i := newItemID(tablename.New(tid), id)

	want := fmt.Sprintf("%v|%v", tid.String(), id)

	if s := i.String(); s != want {
		t.Errorf("incorrect string: %v, expected: %v", i, want)
	}
}
