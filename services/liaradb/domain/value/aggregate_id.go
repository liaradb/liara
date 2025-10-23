package value

type AggregateID struct {
	baseString
}

func NewAggregateID(value string) AggregateID {
	return AggregateID{baseString(value)}
}
