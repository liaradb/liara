package tablename

import (
	"fmt"

	"github.com/liaradb/liaradb/domain/value"
)

type TableName struct {
	n string
}

func New(n string) TableName {
	return TableName{n}
}

func (tn *TableName) KeyValue(pid value.PartitionID) string {
	return fmt.Sprintf("%v--%v.kv", tn.n, pid)
}

func (tn *TableName) EventLog(pid value.PartitionID) string {
	return fmt.Sprintf("%v--%v.el", tn.n, pid)
}

func (tn *TableName) Index(i int, pid value.PartitionID) string {
	return fmt.Sprintf("%v--%v--%v.idx", tn.n, i, pid)
}
