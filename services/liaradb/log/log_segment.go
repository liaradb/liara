package log

import (
	"io/fs"
)

type LogSegment struct {
}

func GetSegmentForLSN(names []LogSegmentName, lsn LogSequenceNumber) (LogSegmentName, bool) {
	for i := len(names) - 1; i >= 0; i-- {
		n := names[i]
		if lsn >= n.logSequenceNumber {
			return n, true
		}
	}

	return LogSegmentName{}, false
}

func ListSegments(fsys fs.FS, dir string) ([]LogSegmentName, error) {
	files, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}

	names := make([]LogSegmentName, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, ParseLogSegmentName(f.Name()))
		}
	}

	return names, nil
}
