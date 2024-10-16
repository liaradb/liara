package liara

type OutboxID string

func (o OutboxID) String() string { return string(o) }
