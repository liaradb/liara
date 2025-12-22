package tablename

import (
	"fmt"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/storage/link"
)

type TableName struct {
	n string
}

func New(n string) TableName {
	return TableName{n}
}

func (tn *TableName) KeyValue(pid value.PartitionID) link.FileName {
	return link.NewFileName(fmt.Sprintf("%v--%v.kv", tn.n, pid))
}

func (tn *TableName) EventLog(pid value.PartitionID) link.FileName {
	return link.NewFileName(fmt.Sprintf("%v--%v.el", tn.n, pid))
}

func (tn *TableName) Index(i int, pid value.PartitionID) link.FileName {
	return link.NewFileName(fmt.Sprintf("%v--%v--%v.idx", tn.n, i, pid))
}
