package schema

type ByteColumn struct {
	name string
}

func NewByteColumn(name string) ByteColumn {
	return ByteColumn{
		name: name,
	}
}

func (ic ByteColumn) Name() string     { return ic.name }
func (ic ByteColumn) Size() int        { return 1 }
func (ic ByteColumn) Type() ColumnType { return ColumnTypeByte }

// TODO: How do we read this value?
func (ic ByteColumn) Value() byte {
	return 0
}
