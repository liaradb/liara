package value

type PartitionID int32

func (p PartitionID) Value() int32 { return int32(p) }
