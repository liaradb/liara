package log

import (
	"io/fs"
)

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

func GetSegmentForLSN(names []SegmentName, lsn LogSequenceNumber) (SegmentName, bool) {
	for i := len(names) - 1; i >= 0; i-- {
		n := names[i]
		if lsn >= n.lsn {
			return n, true
		}
	}

	return SegmentName{}, false
}

func ListSegments(fsys fs.FS, dir string) ([]SegmentName, error) {
	files, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}

	names := make([]SegmentName, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, ParseSegmentName(f.Name()))
		}
	}

	return names, nil
}
