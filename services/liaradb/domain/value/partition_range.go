package value

type PartitionRange struct {
	low  PartitionID
	high PartitionID
}

func NewPartitionRange(pids ...PartitionID) PartitionRange {
	var low PartitionID
	var high PartitionID

	switch len(pids) {
	case 0:
		break
	case 1:
		low, high = pids[0], pids[0]
	default:
		low, high = pids[0], pids[1]
	}

	if low > high {
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
