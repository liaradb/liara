package log

import (
	"io/fs"
)

type Segment struct {
}

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
			names = append(names, ParseLogSegmentName(f.Name()))
		}
	}

	return names, nil
}
