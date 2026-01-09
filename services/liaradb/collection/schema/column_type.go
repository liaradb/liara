package schema

type ColumnType int

const (
	ColumnTypeUnspecified ColumnType = iota
	ColumnTypeBool
	ColumnTypeByte
	ColumnTypeInt8
	ColumnTypeInt16
	ColumnTypeInt32
	ColumnTypeInt64
	ColumnTypeUInt8
	ColumnTypeUInt16
	ColumnTypeUInt32
	ColumnTypeUInt64
	ColumnTypeString
)
