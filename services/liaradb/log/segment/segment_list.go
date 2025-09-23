package segment

import (
	"container/list"
	"io/fs"
	"iter"
	"slices"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/record"
)

type SegmentList struct {
	dir   string
	file  file.File
	fsys  file.FileSystem
	names *list.List
}

func NewSegmentList(fsys file.FileSystem, dir string) *SegmentList {
	return &SegmentList{
		dir:  dir,
		fsys: fsys,
	}
}

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

func (sl *SegmentList) Names() []SegmentName {
	names := make([]SegmentName, 0, sl.names.Len())
	for sf := range sl.iterate() {
		names = append(names, sf.SegmentName())
	}
	return names
}

// TODO: Should not open files unless this is called
func (sl *SegmentList) Open() error {
	files, err := sl.getFiles()
	if err != nil {
		return err
	}

	sl.names = sl.filesToNames(files)
	return nil
}

func (sl *SegmentList) OpenLatestSegment() (SegmentName, file.File, error) {
	if err := sl.Close(); err != nil {
		return SegmentName{}, nil, err
	}

	var sn SegmentName
	sf, ok := sl.getLatestSegment()
	if ok {
		// first segment file
		sn = sf.SegmentName()
	}
	f, err := sl.fsys.OpenFile(sn.String())
	if err != nil {
		return SegmentName{}, nil, err
	}

	sl.file = f
	if !ok {
		sl.names.PushBack(newSegmentFile(sn, sl.fsys))
	}

	return sn, f, err
}

func (sl *SegmentList) OpenNextSegment(lsn record.LogSequenceNumber) (SegmentName, file.File, error) {
	if err := sl.Close(); err != nil {
		return SegmentName{}, nil, err
	}

	var sn SegmentName
	sf, ok := sl.getNextSegment(lsn)
	if ok {
		sn = sf.SegmentName()
	} else {
		sf = newSegmentFile(sn, sl.fsys)
	}
	f, err := sl.fsys.OpenFile(sn.String())
	if err != nil {
		return SegmentName{}, nil, err
	}

	sl.file = f
	sl.names.PushBack(sf)

	return sn, f, err
}

func (sl *SegmentList) OpenSegmentBeforeLSN(lsn record.LogSequenceNumber) (SegmentName, file.File, error) {
	if err := sl.Close(); err != nil {
		return SegmentName{}, nil, err
	}

	sf, _, ok := sl.getSegmentBeforeLSN(lsn)
	if !ok {
		return SegmentName{}, nil, ErrNoSegmentFile
	}

	sn := sf.SegmentName()
	f, err := sl.fsys.OpenFile(sn.String())
	if err != nil {
		return SegmentName{}, nil, err
	}

	sl.file = f
	return sn, f, nil
}

func (sl *SegmentList) OpenSegmentForLSN(lsn record.LogSequenceNumber) (SegmentName, file.File, error) {
	if err := sl.Close(); err != nil {
		return SegmentName{}, nil, err
	}

	sf, ok := sl.getSegmentForLSN(lsn)
	if !ok {
		return SegmentName{}, nil, ErrNoSegmentFile
	}

	sn := sf.SegmentName()
	f, err := sl.fsys.OpenFile(sn.String())
	if err != nil {
		return SegmentName{}, nil, err
	}

	sl.file = f

	return sn, f, err
}

// TODO: Handle if this file is open
func (sl *SegmentList) RemoveSegmentBeforeLSN(lsn record.LogSequenceNumber) error {
	sf, e, ok := sl.getSegmentBeforeLSN(lsn)
	if !ok {
		return ErrNoSegmentFile
	}

	sn := sf.SegmentName()
	if err := sl.fsys.Remove(sn.String()); err != nil {
		return err
	}

	sl.names.Remove(e)
	return nil
}

func (sl *SegmentList) filesToNames(files []fs.DirEntry) *list.List {
	names := make([]*segmentFile, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, newSegmentFile(
				ParseSegmentName(f.Name()),
				sl.fsys))
		}
	}
	// TODO: Do we need to sort?
	slices.SortFunc(names, func(a, b *segmentFile) int {
		return int(a.sn.ID() - b.sn.ID())
	})

	l := list.New()
	for _, sf := range names {
		l.PushBack(sf)
	}

	return l
}

func (sl *SegmentList) getFiles() ([]fs.DirEntry, error) {
	return sl.fsys.ReadDir(sl.dir)
}

func (sl *SegmentList) getLatestSegment() (*segmentFile, bool) {
	e := sl.names.Back()
	if e == nil {
		return nil, false
	}

	return e.Value.(*segmentFile), true
}

func (sl *SegmentList) getNextSegment(lsn record.LogSequenceNumber) (*segmentFile, bool) {
	sf, ok := sl.getLatestSegment()
	if !ok {
		return nil, false
	}

	return sf.Next(lsn), true
}

func (sl *SegmentList) getSegmentBeforeLSN(lsn record.LogSequenceNumber) (*segmentFile, *list.Element, bool) {
	for sf, e := range sl.iterate() {
		if lsn >= sf.SegmentName().LogSequenceNumber() {
			return sf, e, true
		}
	}

	return nil, nil, false
}

func (sl *SegmentList) getSegmentForLSN(lsn record.LogSequenceNumber) (*segmentFile, bool) {
	for sf := range sl.reverse() {
		if lsn >= sf.SegmentName().LogSequenceNumber() {
			return sf, true
		}
	}

	return nil, false
}

func (sl *SegmentList) iterate() iter.Seq2[*segmentFile, *list.Element] {
	return func(yield func(*segmentFile, *list.Element) bool) {
		if sl.names == nil {
			return
		}

		e := sl.names.Front()
		for {
			if e == nil {
				return
			}

			if !yield(e.Value.(*segmentFile), e) {
				return
			}

			e = e.Next()
		}
	}
}

func (sl *SegmentList) reverse() iter.Seq[*segmentFile] {
	return func(yield func(*segmentFile) bool) {
		if sl.names == nil {
			return
		}

		e := sl.names.Back()
		for {
			if e == nil {
				return
			}

			if !yield(e.Value.(*segmentFile)) {
				return
			}

			e = e.Prev()
		}
	}
}
