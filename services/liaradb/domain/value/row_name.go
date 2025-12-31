package value

type RowName struct {
	baseString
}

func NewRowName(value string) RowName {
	return RowName{baseString(value)}
}
