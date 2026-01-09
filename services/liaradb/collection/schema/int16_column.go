package schema

type Int16Column struct {
	name string
}

func NewInt16Column(name string) Int16Column {
	return Int16Column{
		name: name,
	}
}

func (ic Int16Column) Name() string     { return ic.name }
func (ic Int16Column) Size() int        { return 2 }
func (ic Int16Column) Type() ColumnType { return ColumnTypeInt16 }
