package tablename

import (
	"fmt"
	"path"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/storage/link"
)

type Table struct {
	tf   TenantFolder
	name string
}

func newTable(tf TenantFolder, name string) Table {
	return Table{
		tf:   tf,
		name: name,
	}
}

func (t Table) Collection() link.FileName {
	return link.NewFileName(path.Join(t.tf.String(), fmt.Sprintf("%v.col", t.name)))
}

func (t Table) Index(i base.Uint32, pid value.PartitionID) link.FileName {
	return link.NewFileName(path.Join(t.tf.String(), fmt.Sprintf("%v--%v--%v.bt", t.name, i, pid)))
}
