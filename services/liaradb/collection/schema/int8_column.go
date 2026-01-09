package schema

type Int8Column struct {
	name string
}

func NewInt8Column(name string) Int8Column {
	return Int8Column{
		name: name,
	}
}

func (ic Int8Column) Name() string     { return ic.name }
func (ic Int8Column) Size() int        { return 1 }
func (ic Int8Column) Type() ColumnType { return ColumnTypeInt8 }
