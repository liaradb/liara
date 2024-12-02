package value

type RequestID string

func (i RequestID) String() string { return string(i) }
