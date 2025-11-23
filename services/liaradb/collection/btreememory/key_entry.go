package btreememory

import (
	"cmp"
	"fmt"

	"github.com/liaradb/liaradb/storage"
)

type keyEntry[K cmp.Ordered] struct {
	k  K
	id storage.Offset
}

func (ke keyEntry[K]) String() string {
	return fmt.Sprintf("(%v -> %v)", ke.k, ke.id)
}
