package value

type Schema struct {
	baseString
}

func NewSchema(value string) Schema {
	return Schema{baseString(value)}
}
