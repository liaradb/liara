package value

type IdempotenceID string

func (i IdempotenceID) String() string { return string(i) }
