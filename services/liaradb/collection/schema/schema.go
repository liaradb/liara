package schema

type Schema struct {
	id      ID
	columns []Column
}

func (s *Schema) ID() ID { return s.id }

type Column interface {
	Name() string
	Size() int
	Type() ColumnType
}
