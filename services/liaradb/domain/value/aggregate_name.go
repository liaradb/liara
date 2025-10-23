package value

type AggregateName struct {
	baseString
}

func NewAggregateName(value string) AggregateName {
	return AggregateName{baseString(value)}
}
