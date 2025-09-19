package log

import (
	"io/fs"
	"slices"

	"github.com/liaradb/liaradb/file"
)

// TODO: Remove files before LSN

type SegmentList struct {
	fsys  file.FileSystem
	dir   string
	names []SegmentName
}

type ReadDir interface {
	ReadDir(name string) ([]fs.DirEntry, error)
}

func NewSegmentList(fsys file.FileSystem, dir string) *SegmentList {
	return &SegmentList{
		fsys: fsys,
		dir:  dir,
	}
}

func (sl *SegmentList) Names() []SegmentName { return sl.names }

func (sl *SegmentList) Open() error {
	files, err := sl.getFiles()
	if err != nil {
		return err
	}

	names := sl.filesToNames(files)
	sl.names = names
	return nil
}

func (sl *SegmentList) getFiles() ([]fs.DirEntry, error) {
	return sl.fsys.ReadDir(sl.dir)
}

func (*SegmentList) filesToNames(files []fs.DirEntry) []SegmentName {
	names := make([]SegmentName, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, ParseSegmentName(f.Name()))
		}
	}
	slices.SortFunc(names, func(a, b SegmentName) int {
		return int(a.ID()) - int(b.ID())
	})
	return names
}

func (sl *SegmentList) OpenLatestSegment() (SegmentName, file.File, error) {
	sn, ok := sl.getLatestSegment()
	f, err := sl.fsys.OpenFile(sn.String())
	if err != nil {
		return SegmentName{}, nil, err
	}

	if !ok {
		sl.names = append(sl.names, sn)
	}

	return sn, f, err
}

func (sl *SegmentList) getLatestSegment() (SegmentName, bool) {
	if len(sl.names) > 0 {
		return sl.names[len(sl.names)-1], true
	}

	return SegmentName{}, false
}

func (sl *SegmentList) OpenSegmentForLSN(lsn LogSequenceNumber) (SegmentName, file.File, error) {
	sn, ok := sl.getSegmentForLSN(lsn)
	if !ok {
		return SegmentName{}, nil, ErrNoSegmentFile
	}
	f, err := sl.fsys.OpenFile(sn.String())
	return sn, f, err
}

func (sl *SegmentList) getSegmentForLSN(lsn LogSequenceNumber) (SegmentName, bool) {
	for i := len(sl.names) - 1; i >= 0; i-- {
		n := sl.names[i]
		if lsn >= n.lsn {
			return n, true
		}
	}

	return SegmentName{}, false
}
