package liara

type Version int

func (v Version) Value() int { return int(v) }
