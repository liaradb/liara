package eventsource

type Data []byte

func (d Data) String() string { return string(d) }
