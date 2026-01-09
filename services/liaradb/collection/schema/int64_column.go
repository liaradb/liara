package schema

type Int64Column struct {
	name string
}

func (ic Int64Column) Name() string {
	return ic.name
}

func (ic Int64Column) Size() int {
	return 8
}

func (ic Int64Column) Type() ColumnType {
	return ColumnTypeInt64
}
