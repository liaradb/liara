package tablename

import (
	"fmt"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/storage/link"
)

const defaultTenantID = "default"

type TableName struct {
	value string
}

func New(tenantID value.TenantID) TableName {
	return TableName{
		value: tenantID.String(),
	}
}

func NewFromString(value string) TableName {
	return TableName{
		value: value,
	}
}

func (tn *TableName) String() string {
	if tn.value == "" {
		return defaultTenantID
	}

	return string(tn.value)
}

func (tn *TableName) KeyValue(pid value.PartitionID) link.FileName {
	return link.NewFileName(fmt.Sprintf("%v--%v.kv", tn, pid))
}

func (tn *TableName) EventLog(pid value.PartitionID) link.FileName {
	return link.NewFileName(fmt.Sprintf("%v--%v.el", tn, pid))
}

func (tn *TableName) Outbox(pid value.PartitionID) link.FileName {
	return link.NewFileName(fmt.Sprintf("%v--%v.out", tn, pid))
}

func (tn *TableName) RequestLog() link.FileName {
	return link.NewFileName(fmt.Sprintf("%v.rl", tn))
}

func (tn *TableName) Index(i base.Uint32, pid value.PartitionID) link.FileName {
	return link.NewFileName(fmt.Sprintf("%v--%v--%v.idx", tn, i, pid))
}
