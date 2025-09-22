package segment

import (
	"container/list"
	"fmt"
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
	for n := range sl.iterate() {
		names = append(names, n)
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

	sn, ok := sl.getLatestSegment()
	f, err := sl.fsys.OpenFile(sn.String())
	if err != nil {
		return SegmentName{}, nil, err
	}

	sl.file = f
	if !ok {
		sl.names.PushBack(sn)
	}

	return sn, f, err
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
	sl.names.PushBack(sn)

	return sn, f, err
}

func (sl *SegmentList) OpenSegmentBeforeLSN(lsn record.LogSequenceNumber) (SegmentName, file.File, error) {
	if err := sl.Close(); err != nil {
		return SegmentName{}, nil, err
	}

	sn, _, ok := sl.getSegmentBeforeLSN(lsn)
	if !ok {
		return SegmentName{}, nil, ErrNoSegmentFile
	}

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

// TODO: Handle if this file is open
func (sl *SegmentList) RemoveSegmentBeforeLSN(lsn record.LogSequenceNumber) error {
	sn, e, ok := sl.getSegmentBeforeLSN(lsn)
	if !ok {
		return ErrNoSegmentFile
	}

	if err := sl.fsys.Remove(sn.String()); err != nil {
		return err
	}

	sl.names.Remove(e)
	return nil
}

func (*SegmentList) filesToNames(files []fs.DirEntry) *list.List {
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

	l := list.New()
	for _, n := range names {
		l.PushBack(n)
	}

	return l
}

func (sl *SegmentList) getFiles() ([]fs.DirEntry, error) {
	return sl.fsys.ReadDir(sl.dir)
}

func (sl *SegmentList) getLatestSegment() (SegmentName, bool) {
	e := sl.names.Back()
	if e == nil {
		return SegmentName{}, false
	}

	return e.Value.(SegmentName), true
}

func (sl *SegmentList) getNextSegment(lsn record.LogSequenceNumber) SegmentName {
	sn, ok := sl.getLatestSegment()
	if !ok {
		return SegmentName{}
	}

	return sn.Next(lsn)
}

func (sl *SegmentList) getSegmentBeforeLSN(lsn record.LogSequenceNumber) (SegmentName, *list.Element, bool) {
	for n, e := range sl.iterate() {
		if lsn >= n.lsn {
			return n, e, true
		}
	}

	return SegmentName{}, nil, false
}

func (sl *SegmentList) getSegmentForLSN(lsn record.LogSequenceNumber) (SegmentName, bool) {
	for n := range sl.reverse() {
		fmt.Printf("%v, %v, %v\n", lsn, n.lsn, lsn >= n.lsn)
		if lsn >= n.lsn {
			return n, true
		}
	}

	return SegmentName{}, false
}

func (sl *SegmentList) iterate() iter.Seq2[SegmentName, *list.Element] {
	return func(yield func(SegmentName, *list.Element) bool) {
		if sl.names == nil {
			return
		}

		e := sl.names.Front()
		for {
			if e == nil {
				return
			}

			if !yield(e.Value.(SegmentName), e) {
				return
			}

			e = e.Next()
		}
	}
}

func (sl *SegmentList) reverse() iter.Seq[SegmentName] {
	return func(yield func(SegmentName) bool) {
		if sl.names == nil {
			return
		}

		e := sl.names.Back()
		for {
			if e == nil {
				return
			}

			if !yield(e.Value.(SegmentName)) {
				return
			}

			e = e.Prev()
		}
	}
}
