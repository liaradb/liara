package schema

type Schema struct {
	columns []Column
}

type Column interface {
	Name() string
	Size() int
	Type() ColumnType
}
