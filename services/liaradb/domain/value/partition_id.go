package value

import "github.com/liaradb/liaradb/encoder/base"

type PartitionID struct {
	baseUint32
}

func NewPartitionID(value int32) PartitionID {
	return PartitionID{baseUint32(value)}
}

const PartitionIDSize = base.BaseUint32Size
