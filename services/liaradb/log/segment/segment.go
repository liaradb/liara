package segment

type Segment struct {
	size     int // number of pages
	pageSize int // page size
}

func NewSegment(size int, pageSize int) *Segment {
	return &Segment{
		size:     size,
		pageSize: pageSize,
	}
}

func (s *Segment) Size() int     { return s.size }
func (s *Segment) PageSize() int { return s.pageSize }
