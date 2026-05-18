package segment

type SegmentID uint64

func (id SegmentID) Next() SegmentID {
	return id + 1
}
