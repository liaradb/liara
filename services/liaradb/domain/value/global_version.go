package value

type GlobalVersion int

func (gv GlobalVersion) Value() int { return int(gv) }
