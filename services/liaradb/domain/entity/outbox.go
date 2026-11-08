package entity

import (
	"io"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/serializer"
)

const (
	OutboxSize = value.OutboxIDSize +
		value.PartitionIDSize +
		value.GlobalVersionSize
)

type Outbox struct {
	id            value.OutboxID
	partitionID   value.PartitionID
	globalVersion value.GlobalVersion
}

func NewOutbox(
	id value.OutboxID,
	partitionID value.PartitionID,
) *Outbox {
	return &Outbox{
		id:          id,
		partitionID: partitionID,
	}
}

func RestoreOutbox(
	id value.OutboxID,
	partitionRange value.PartitionID,
	globalVersion value.GlobalVersion,
) *Outbox {
	return &Outbox{
		id:            id,
		partitionID:   partitionRange,
		globalVersion: globalVersion,
	}
}

func (o *Outbox) ID() value.OutboxID                 { return o.id }
func (o *Outbox) PartitionID() value.PartitionID     { return o.partitionID }
func (o *Outbox) GlobalVersion() value.GlobalVersion { return o.globalVersion }
func (o *Outbox) Size() int                          { return OutboxSize }

func (o *Outbox) UpdateGlobalVersion(v value.GlobalVersion) {
	o.globalVersion = v
}

func (o *Outbox) WriteData(data []byte) []byte {
	data0 := o.id.WriteData(data)
	data1 := o.partitionID.WriteData(data0)
	return o.globalVersion.WriteData(data1)
}

func (o *Outbox) ReadData(data []byte) []byte {
	data0 := o.id.ReadData(data)
	data1 := o.partitionID.ReadData(data0)
	return o.globalVersion.ReadData(data1)
}

func (o *Outbox) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&o.id,
		&o.partitionID,
		&o.globalVersion)
}

func (o *Outbox) Write(r io.Writer) error {
	return serializer.WriteAll(r,
		o.id,
		o.partitionID,
		o.globalVersion)
}
