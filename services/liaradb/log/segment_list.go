package log

import (
	"io/fs"
	"slices"

	"github.com/liaradb/liaradb/file"
)

type SegmentList struct {
	fsys file.FileSystem
	dir  string
}

type ReadDir interface {
	ReadDir(name string) ([]fs.DirEntry, error)
}

func NewSegmentList(
	fsys file.FileSystem,
	dir string,
) *SegmentList {
	return &SegmentList{
		fsys: fsys,
		dir:  dir,
	}
}

func (sl *SegmentList) GetLatestSegment(names []SegmentName) SegmentName {
	if len(names) > 0 {
		return names[len(names)-1]
	}

	return SegmentName{}
}

func (sl *SegmentList) GetSegmentForLSN(names []SegmentName, lsn LogSequenceNumber) (SegmentName, bool) {
	for i := len(names) - 1; i >= 0; i-- {
		n := names[i]
		if lsn >= n.lsn {
			return n, true
		}
	}

	return SegmentName{}, false
}

func (sl *SegmentList) ListSegments(fsys ReadDir, dir string) ([]SegmentName, error) {
	files, err := fsys.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	names := make([]SegmentName, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, ParseSegmentName(f.Name()))
		}
	}

	slices.SortFunc(names, func(a, b SegmentName) int {
		return int(a.ID()) - int(b.ID())
	})

	return names, nil
}

func (sl *SegmentList) OpenLatestSegment(fsys file.FileSystem, dir string) (file.File, error) {
	names, err := sl.ListSegments(fsys, ".")
	if err != nil {
		return nil, err
	}

	sn := sl.GetLatestSegment(names)
	return fsys.Open(sn.String())
}
