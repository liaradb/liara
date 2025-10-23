package value

import "github.com/liaradb/liaradb/raw"

type AggregateID struct {
	baseString
}

func NewAggregateID() AggregateID {
	return AggregateID{raw.NewBaseString()}
}
