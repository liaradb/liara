package liara

type Limit int

func (l Limit) Value() int { return int(l) }
