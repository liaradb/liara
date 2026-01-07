package tablename

import (
	"fmt"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/storage/link"
)

type TableName struct {
	tenantID value.TenantID
}

func New(tenantID value.TenantID) TableName {
	return TableName{
		tenantID: tenantID,
	}
}

func (tn *TableName) KeyValue(pid value.PartitionID) link.FileName {
	return link.NewFileName(fmt.Sprintf("%v--%v.kv", tn.tenantID, pid))
}

func (tn *TableName) EventLog(pid value.PartitionID) link.FileName {
	return link.NewFileName(fmt.Sprintf("%v--%v.el", tn.tenantID, pid))
}

func (tn *TableName) Index(i raw.BaseUint32, pid value.PartitionID) link.FileName {
	return link.NewFileName(fmt.Sprintf("%v--%v--%v.idx", tn.tenantID, i, pid))
}
