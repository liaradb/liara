package tablename

import (
	"fmt"
	"path"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/storage/link"
)

type PartitionedTableName struct {
	tf   TenantFolder
	name string
	pid  value.PartitionID
}

func newPartitionedTable(tf TenantFolder, name string, pid value.PartitionID) PartitionedTableName {
	return PartitionedTableName{
		tf:   tf,
		name: name,
		pid:  pid,
	}
}

func (t PartitionedTableName) Collection() link.FileName {
	return link.NewFileName(path.Join(t.tf.String(), fmt.Sprintf("%v--%v.col", t.name, t.pid)))
}

func (t PartitionedTableName) Index(i base.Uint32) link.FileName {
	return link.NewFileName(path.Join(t.tf.String(), fmt.Sprintf("%v--%v--%v.bt", t.name, i, t.pid)))
}
