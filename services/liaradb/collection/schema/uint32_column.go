package schema

type UInt32Column struct {
	name string
}

func NewUInt32Column(name string) UInt32Column {
	return UInt32Column{
		name: name,
	}
}

func (ic UInt32Column) Name() string     { return ic.name }
func (ic UInt32Column) Size() int        { return 4 }
func (ic UInt32Column) Type() ColumnType { return ColumnTypeUInt32 }

func (ic UInt32Column) Value() uint32 {
	return 0
}
