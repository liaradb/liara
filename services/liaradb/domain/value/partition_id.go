package value

type PartitionID struct {
	baseUint32
}

func NewPartitionID(value uint32) PartitionID {
	return PartitionID{baseUint32(value)}
}
