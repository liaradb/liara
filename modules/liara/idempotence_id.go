package liara

type IdempotenceID string

func (i IdempotenceID) String() string { return string(i) }
