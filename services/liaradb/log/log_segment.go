package log

import (
	"os"
)

type LogSegment struct {
}

func ListSegments(dir string) ([]LogSegmentName, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	names := make([]LogSegmentName, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, LogSegmentName(f.Name()))
		}
	}
	return names, nil
}
