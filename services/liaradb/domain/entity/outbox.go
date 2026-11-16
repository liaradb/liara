package entity

import (
	"io"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/serializer"
)

const (
	OutboxSize = value.OutboxIDSize +
		value.PartitionIDSize +
		value.GlobalVersionSize +
		value.SubVersionSize
)

type Outbox struct {
	id            value.OutboxID
	partitionID   value.PartitionID
	globalVersion value.GlobalVersion
	subVersion    value.SubVersion
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
	subVersion value.SubVersion,
) *Outbox {
	return &Outbox{
		id:            id,
		partitionID:   partitionRange,
		globalVersion: globalVersion,
		subVersion:    subVersion,
	}
}

func (o *Outbox) ID() value.OutboxID                 { return o.id }
func (o *Outbox) PartitionID() value.PartitionID     { return o.partitionID }
func (o *Outbox) GlobalVersion() value.GlobalVersion { return o.globalVersion }
func (o *Outbox) SubVersion() value.SubVersion       { return o.subVersion }
func (o *Outbox) Size() int                          { return OutboxSize }

func (o *Outbox) UpdatePosition(v value.GlobalVersion, sv value.SubVersion) {
	o.globalVersion = v
	o.subVersion = sv
}

func (o *Outbox) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&o.id,
		&o.partitionID,
		&o.globalVersion,
		&o.subVersion)
}

func (o *Outbox) Write(r io.Writer) error {
	return serializer.WriteAll(r,
		o.id,
		o.partitionID,
		o.globalVersion,
		o.subVersion)
}
