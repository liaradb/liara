package entity

import "github.com/liaradb/liaradb/domain/value"

const (
	OutboxSize = value.OutboxIDSize +
		value.PartitionRangeSize +
		value.GlobalVersionSize
)

type Outbox struct {
	id             value.OutboxID
	partitionRange value.PartitionRange
	globalVersion  value.GlobalVersion
}

func NewOutbox(
	id value.OutboxID,
	partitionRange value.PartitionRange,
) *Outbox {
	return &Outbox{
		id:             id,
		partitionRange: partitionRange,
	}
}

func RestoreOutbox(
	id value.OutboxID,
	partitionRange value.PartitionRange,
	globalVersion value.GlobalVersion,
) *Outbox {
	return &Outbox{
		id:             id,
		partitionRange: partitionRange,
		globalVersion:  globalVersion,
	}
}

func (o *Outbox) ID() value.OutboxID                   { return o.id }
func (o *Outbox) PartitionRange() value.PartitionRange { return o.partitionRange }
func (o *Outbox) GlobalVersion() value.GlobalVersion   { return o.globalVersion }

func (o *Outbox) UpdateGlobalVersion(v value.GlobalVersion) {
	o.globalVersion = v
}

func (o *Outbox) Write(data []byte) ([]byte, bool) {
	data0, ok := o.globalVersion.WriteData(data)
	if !ok {
		return nil, false
	}

	data1 := o.partitionRange.WriteData(data0)
	return o.id.WriteData(data1), true
}

func (o *Outbox) Read(data []byte) ([]byte, bool) {
	data0, ok := o.globalVersion.ReadData(data)
	if !ok {
		return nil, false
	}

	data1 := o.partitionRange.ReadData(data0)
	return o.id.ReadData(data1), true
}
