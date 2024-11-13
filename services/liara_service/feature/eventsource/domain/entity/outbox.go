package entity

import "github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"

type Outbox struct {
	id            value.OutboxID
	aggregateName value.AggregateName
	globalVersion value.GlobalVersion
}

func (o *Outbox) ID() value.OutboxID                 { return o.id }
func (o *Outbox) AggregateName() value.AggregateName { return o.aggregateName }
func (o *Outbox) GlobalVersion() value.GlobalVersion { return o.globalVersion }

func NewOutbox(
	id value.OutboxID,
	aggregateName value.AggregateName,
) *Outbox {
	return &Outbox{
		id:            id,
		aggregateName: aggregateName,
	}
}
