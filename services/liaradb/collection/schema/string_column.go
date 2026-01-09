package schema

type StringColumn struct {
	name string
	size int
}

func NewStringColumn(name string, size int) StringColumn {
	return StringColumn{
		name: name,
		size: size,
	}
}

func (sc StringColumn) Name() string     { return sc.name }
func (sc StringColumn) Size() int        { return sc.size }
func (sc StringColumn) Type() ColumnType { return ColumnTypeString }
