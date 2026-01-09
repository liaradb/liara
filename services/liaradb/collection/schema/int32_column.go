package schema

type Int32Column struct {
	name string
}

func NewInt32Column(name string) Int32Column {
	return Int32Column{
		name: name,
	}
}

func (ic Int32Column) Name() string     { return ic.name }
func (ic Int32Column) Size() int        { return 4 }
func (ic Int32Column) Type() ColumnType { return ColumnTypeInt32 }
