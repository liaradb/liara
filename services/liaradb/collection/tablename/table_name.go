package tablename

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/storage/link"
)

var defaultTenantID = value.NewTenantIDFromUUID(uuid.Nil)

type TableName struct {
	tenantID value.TenantID
	name     string
}

func New(tenantID value.TenantID) TableName {
	return TableName{
		tenantID: tenantID,
		name:     "",
	}
}

// TODO: Should this be used?
func NewFromString(name string) TableName {
	return TableName{
		tenantID: defaultTenantID,
		name:     name,
	}
}

func (tn *TableName) TenantID() value.TenantID { return tn.tenantID }

func (tn *TableName) String() string {
	if tn.name == "" {
		return tn.tenantID.String()
	}

	return fmt.Sprintf("%v--%v", tn.tenantID, tn.name)
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
