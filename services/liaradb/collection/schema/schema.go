package schema

type Schema struct {
	id      ID
	columns []Column
	names   map[string]Column
}

type Column interface {
	Name() string
	Size() int
	Type() ColumnType
}

func NewSchema(id ID) *Schema {
	return &Schema{
		id: id,
	}
}

func (s *Schema) ID() ID { return s.id }

func (s *Schema) AddColumn(c Column) bool {
	name := c.Name()
	if _, ok := s.names[name]; ok {
		return false
	}

	s.columns = append(s.columns, c)
	if s.names == nil {
		s.names = make(map[string]Column)
	}
	s.names[name] = c
	return true
}

func (s *Schema) GetString(name string) (string, bool) {
	c, ok := s.names[name]
	if !ok {
		return "", false
	}

	sc, ok := c.(StringColumn)
	if !ok {
		return "", false
	}

	return sc.Value(), true
}
