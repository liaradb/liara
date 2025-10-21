package value

type EventName string

func (n EventName) String() string { return string(n) }
