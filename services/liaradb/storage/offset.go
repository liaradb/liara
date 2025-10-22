package storage

type Offset int64

func (o Offset) Value() int64 { return int64(o) }
