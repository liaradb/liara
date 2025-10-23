package value

type CorrelationID struct {
	baseString
}

func NewCorrelationID(value string) CorrelationID {
	return CorrelationID{baseString(value)}
}
