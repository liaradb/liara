package tenantfolder

import (
	"fmt"
	"path"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/storage/link"
)

type GlobalTableName struct {
	tf   TenantFolder
	name string
}

func newGlobalTableName(tf TenantFolder, name string) GlobalTableName {
	return GlobalTableName{
		tf:   tf,
		name: name,
	}
}

func (t GlobalTableName) Collection() link.FileName {
	return tablename.CollectionExtension.FileName(t.String())
}

func (t GlobalTableName) Index(i base.Uint32) link.FileName {
	return tablename.BTreeExtension.FileName(path.Join(t.tf.String(), fmt.Sprintf("%v--%v", t.name, i)))
}

func (t GlobalTableName) String() string {
	return path.Join(t.tf.String(), t.name)
}
