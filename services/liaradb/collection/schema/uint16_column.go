package schema

type UInt16Column struct {
	name string
}

func NewUInt16Column(name string) UInt16Column {
	return UInt16Column{
		name: name,
	}
}

func (ic UInt16Column) Name() string     { return ic.name }
func (ic UInt16Column) Size() int        { return 2 }
func (ic UInt16Column) Type() ColumnType { return ColumnTypeUInt16 }

// TODO: How do we read this value?
func (ic UInt16Column) Value() uint16 {
	return 0
}
