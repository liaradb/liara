package log

import (
	"io/fs"
	"slices"

	"github.com/liaradb/liaradb/file"
)

type SegmentList struct {
	fsys  file.FileSystem
	dir   string
	names []SegmentName
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

func (sl *SegmentList) Names() []SegmentName { return sl.names }

func (sl *SegmentList) Open() error {
	files, err := sl.fsys.ReadDir(sl.dir)
	if err != nil {
		return err
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

	sl.names = names

	return nil
}

func (sl *SegmentList) GetLatestSegment() SegmentName {
	if len(sl.names) > 0 {
		return sl.names[len(sl.names)-1]
	}

	return SegmentName{}
}

func (sl *SegmentList) GetSegmentForLSN(lsn LogSequenceNumber) (SegmentName, bool) {
	for i := len(sl.names) - 1; i >= 0; i-- {
		n := sl.names[i]
		if lsn >= n.lsn {
			return n, true
		}
	}

	return SegmentName{}, false
}

// TODO: Test this
func (sl *SegmentList) OpenLatestSegment() (file.File, error) {
	return sl.fsys.OpenFile(sl.GetLatestSegment().String())
}
