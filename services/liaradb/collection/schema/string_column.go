package schema

type StringColumn struct {
	name string
	size int
}

func (sc StringColumn) Name() string {
	return sc.name
}

func (sc StringColumn) Size() int {
	return sc.size
}

func (sc StringColumn) Type() ColumnType {
	return ColumnTypeString
}
