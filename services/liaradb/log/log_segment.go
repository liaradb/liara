package log

import (
	"io/fs"
)

type LogSegment struct {
}

func ListSegments(dir string, fsys fs.FS) ([]LogSegmentName, error) {
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
