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

func (s *Schema) GetInt8(name string) (int8, bool) {
	c, ok := s.names[name]
	if !ok {
		return 0, false
	}

	ic, ok := c.(Int8Column)
	if !ok {
		return 0, false
	}

	return ic.Value(), true
}

func (s *Schema) GetInt16(name string) (int16, bool) {
	c, ok := s.names[name]
	if !ok {
		return 0, false
	}

	ic, ok := c.(Int16Column)
	if !ok {
		return 0, false
	}

	return ic.Value(), true
}

func (s *Schema) GetInt32(name string) (int32, bool) {
	c, ok := s.names[name]
	if !ok {
		return 0, false
	}

	ic, ok := c.(Int32Column)
	if !ok {
		return 0, false
	}

	return ic.Value(), true
}

func (s *Schema) GetInt64(name string) (int64, bool) {
	c, ok := s.names[name]
	if !ok {
		return 0, false
	}

	ic, ok := c.(Int64Column)
	if !ok {
		return 0, false
	}

	return ic.Value(), true
}

func (s *Schema) GetByte(name string) (byte, bool) {
	c, ok := s.names[name]
	if !ok {
		return 0, false
	}

	ic, ok := c.(ByteColumn)
	if !ok {
		return 0, false
	}

	return ic.Value(), true
}

func (s *Schema) GetUInt16(name string) (uint16, bool) {
	c, ok := s.names[name]
	if !ok {
		return 0, false
	}

	ic, ok := c.(UInt16Column)
	if !ok {
		return 0, false
	}

	return ic.Value(), true
}

func (s *Schema) GetUInt32(name string) (uint32, bool) {
	c, ok := s.names[name]
	if !ok {
		return 0, false
	}

	ic, ok := c.(UInt32Column)
	if !ok {
		return 0, false
	}

	return ic.Value(), true
}

func (s *Schema) GetUInt64(name string) (uint64, bool) {
	c, ok := s.names[name]
	if !ok {
		return 0, false
	}

	ic, ok := c.(UInt64Column)
	if !ok {
		return 0, false
	}

	return ic.Value(), true
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
