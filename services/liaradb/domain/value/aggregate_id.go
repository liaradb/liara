package value

type AggregateID struct {
	baseID
}

func NewAggregateID() AggregateID {
	return AggregateID{newBaseID()}
}
