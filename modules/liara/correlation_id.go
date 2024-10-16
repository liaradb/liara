package liara

type CorrelationID string

func (c CorrelationID) String() string { return string(c) }
