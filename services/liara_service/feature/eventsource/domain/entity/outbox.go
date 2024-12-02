package entity

import "github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"

type Outbox struct {
	id             value.OutboxID
	partitionRange value.PartitionRange
	globalVersion  value.GlobalVersion
}

func (o *Outbox) ID() value.OutboxID                   { return o.id }
func (o *Outbox) PartitionRange() value.PartitionRange { return o.partitionRange }
func (o *Outbox) GlobalVersion() value.GlobalVersion   { return o.globalVersion }

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
