package entity

import (
	"io"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/serializer"
)

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

func (o *Outbox) Write(data []byte) []byte {
	data0 := o.id.WriteData(data)
	data1 := o.partitionRange.WriteData(data0)
	return o.globalVersion.WriteData(data1)
}

func (o *Outbox) ReadData(data []byte) []byte {
	data0 := o.id.ReadData(data)
	data1 := o.partitionRange.ReadData(data0)
	return o.globalVersion.ReadData(data1)
}

func (o *Outbox) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&o.id,
		&o.partitionRange,
		&o.globalVersion)
}
