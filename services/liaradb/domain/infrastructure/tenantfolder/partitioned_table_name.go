package tenantfolder

import (
	"fmt"
	"path"

	"github.com/liaradb/liaradb/collection/tablename"
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
	return tablename.CollectionExtension.FileName(t.String())
}

func (t PartitionedTableName) Index(i base.Uint32) link.FileName {
	return tablename.BTreeExtension.FileName(path.Join(t.tf.String(), fmt.Sprintf("%v--%v--%v", t.name, i, t.pid)))
}

func (t PartitionedTableName) String() string {
	return path.Join(t.tf.String(), fmt.Sprintf("%v--%v", t.name, t.pid))
}
