package command

import "github.com/liaradb/liaradb/domain/value"

type AppendEvent struct {
	TenantID    value.TenantID
	Options     AppendOptions
	PartitionID value.PartitionID
	Events      []EventOptions
}
