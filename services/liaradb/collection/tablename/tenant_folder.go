package tablename

import "github.com/liaradb/liaradb/domain/value"

type TenantFolder struct {
	tenantID value.TenantID
}

func NewTenantFolder(tid value.TenantID) TenantFolder {
	return TenantFolder{tid}
}

func NewTenantFolderFromString(s string) (TenantFolder, error) {
	tid, err := value.NewTenantIDFromString(s)
	if err != nil {
		return TenantFolder{}, err
	}

	return TenantFolder{tid}, nil
}

func (tf TenantFolder) TenantID() value.TenantID { return tf.tenantID }
func (tf TenantFolder) String() string           { return tf.tenantID.String() }

func (tf TenantFolder) GlobalTable(name string) GlobalTableName {
	return newGlobalTableName(tf, name)
}

func (tf TenantFolder) PartitionedTable(name string, pid value.PartitionID) PartitionedTableName {
	return newPartitionedTable(tf, name, pid)
}

func (tf TenantFolder) EventLog(pid value.PartitionID) PartitionedTableName {
	return tf.PartitionedTable("event_log", pid)
}

func (tf TenantFolder) RequestLog() GlobalTableName {
	return tf.GlobalTable("request_log")
}

func (tf TenantFolder) Outbox() GlobalTableName {
	return tf.GlobalTable("outbox")
}

func (tf TenantFolder) Tenant() GlobalTableName {
	return tf.GlobalTable("tenant")
}
