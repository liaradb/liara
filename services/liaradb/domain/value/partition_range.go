package value

const PartitionRangeSize = PartitionIDSize + PartitionIDSize

type PartitionRange struct {
	low  PartitionID
	high PartitionID
}

// TODO: Should this be two parameters?
func NewPartitionRange(low PartitionID, high PartitionID) PartitionRange {
	if low.Value() > high.Value() {
		low, high = high, low
	}

	return PartitionRange{
		low:  low,
		high: high,
	}
}

func (pr PartitionRange) Low() PartitionID                { return pr.low }
func (pr PartitionRange) High() PartitionID               { return pr.high }
func (pr PartitionRange) All() (PartitionID, PartitionID) { return pr.low, pr.high }

func (pr PartitionRange) WriteData(data []byte) []byte {
	data0 := pr.low.WriteData(data)
	return pr.high.WriteData(data0)
}

func (pr *PartitionRange) ReadData(data []byte) []byte {
	data0 := pr.low.ReadData(data)
	return pr.high.ReadData(data0)
}
