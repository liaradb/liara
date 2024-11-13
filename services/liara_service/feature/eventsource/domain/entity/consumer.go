package entity

import "github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"

type Consumer struct {
	aggregateName value.AggregateName
	globalVersion value.GlobalVersion
}
