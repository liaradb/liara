package action

type ItemID string

func (i ItemID) String() string { return string(i) }
