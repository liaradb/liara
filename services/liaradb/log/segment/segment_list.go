package segment

import (
	"io/fs"
	"slices"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/record"
)

type SegmentList struct {
	fsys  file.FileSystem
	dir   string
	names []SegmentName
	file  file.File
}

func NewSegmentList(fsys file.FileSystem, dir string) *SegmentList {
	return &SegmentList{
		fsys: fsys,
		dir:  dir,
	}
}

func (sl *SegmentList) Names() []SegmentName { return sl.names }

func (sl *SegmentList) Close() error {
	if sl.file == nil {
		return nil
	}

	if err := sl.file.Close(); err != nil {
		return err
	}

	sl.file = nil
	return nil
}

// TODO: Should not open files unless this is called
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
	// TODO: Do we need to sort?
	slices.SortFunc(names, func(a, b SegmentName) int {
		return int(a.ID() - b.ID())
	})
	return names
}

func (sl *SegmentList) OpenNextSegment(lsn record.LogSequenceNumber) (SegmentName, file.File, error) {
	if err := sl.Close(); err != nil {
		return SegmentName{}, nil, err
	}

	sn := sl.getNextSegment(lsn)
	f, err := sl.fsys.OpenFile(sn.String())
	if err != nil {
		return SegmentName{}, nil, err
	}

	sl.file = f
	sl.names = append(sl.names, sn)

	return sn, f, err
}

func (sl *SegmentList) getNextSegment(lsn record.LogSequenceNumber) SegmentName {
	if len(sl.names) > 0 {
		return sl.names[len(sl.names)-1].Next(lsn)
	}

	return SegmentName{}
}

func (sl *SegmentList) OpenLatestSegment() (SegmentName, file.File, error) {
	if err := sl.Close(); err != nil {
		return SegmentName{}, nil, err
	}

	sn, ok := sl.getLatestSegment()
	f, err := sl.fsys.OpenFile(sn.String())
	if err != nil {
		return SegmentName{}, nil, err
	}

	sl.file = f
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

func (sl *SegmentList) OpenSegmentForLSN(lsn record.LogSequenceNumber) (SegmentName, file.File, error) {
	if err := sl.Close(); err != nil {
		return SegmentName{}, nil, err
	}

	sn, ok := sl.getSegmentForLSN(lsn)
	if !ok {
		return SegmentName{}, nil, ErrNoSegmentFile
	}
	f, err := sl.fsys.OpenFile(sn.String())
	if err != nil {
		return SegmentName{}, nil, err
	}

	sl.file = f

	return sn, f, err
}

func (sl *SegmentList) getSegmentForLSN(lsn record.LogSequenceNumber) (SegmentName, bool) {
	for i := len(sl.names) - 1; i >= 0; i-- {
		n := sl.names[i]
		if lsn >= n.lsn {
			return n, true
		}
	}

	return SegmentName{}, false
}

func (sl *SegmentList) RemoveSegmentBeforeLSN(lsn record.LogSequenceNumber) error {
	sn, index, ok := sl.getSegmentBeforeLSN(lsn)
	if !ok {
		return ErrNoSegmentFile
	}

	sl.names = sl.names[index:]

	return sl.fsys.Remove(sn.String())
}

func (sl *SegmentList) getSegmentBeforeLSN(lsn record.LogSequenceNumber) (SegmentName, int, bool) {
	index := sl.getIndexForLSN(lsn)
	if index > 0 {
		return sl.names[index-1], index, true
	}

	return SegmentName{}, 0, false
}

func (sl *SegmentList) getIndexForLSN(lsn record.LogSequenceNumber) int {
	for i := len(sl.names) - 1; i >= 0; i-- {
		n := sl.names[i]
		if lsn >= n.lsn {
			return i
		}
	}

	return 0
}
