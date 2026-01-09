package schema

type UInt64Column struct {
	name string
}

func NewUInt64Column(name string) UInt64Column {
	return UInt64Column{
		name: name,
	}
}

func (ic UInt64Column) Name() string     { return ic.name }
func (ic UInt64Column) Size() int        { return 8 }
func (ic UInt64Column) Type() ColumnType { return ColumnTypeUInt64 }
