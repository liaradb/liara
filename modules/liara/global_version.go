package liara

type GlobalVersion int

func (gv GlobalVersion) Value() int { return int(gv) }
