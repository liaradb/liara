package segment

import (
	"container/list"
	"io/fs"
	"iter"
	"slices"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/record"
)

type List struct {
	dir   string
	fsys  file.FileSystem
	names *list.List
	sf    segmentFile
}

func NewList(fsys file.FileSystem, dir string) *List {
	return &List{
		dir:  dir,
		fsys: fsys,
		sf:   *newSegmentFile(fsys, dir),
	}
}

func (l *List) Close() error {
	if err := l.sf.Close(); err != nil {
		return err
	}

	return nil
}

func (l *List) Names() []SegmentName {
	names := make([]SegmentName, 0, l.names.Len())
	for n := range l.iterate() {
		names = append(names, n)
	}
	return names
}

func (l *List) Open() error {
	files, err := l.getFiles()
	if err != nil {
		return err
	}

	l.names = l.filesToNames(files)
	return nil
}

func (l *List) init() error {
	if l.names != nil {
		return nil
	}

	return l.Open()
}

func (l *List) OpenLatestSegment() (SegmentName, file.File, error) {
	if err := l.init(); err != nil {
		return SegmentName{}, nil, err
	}

	sn, ok := l.getLatestSegment()
	f, err := l.sf.Open(sn)
	if err != nil {
		return SegmentName{}, nil, err
	}

	if !ok {
		l.names.PushBack(sn)
	}

	return sn, f, err
}

func (l *List) OpenNextSegment(lsn record.LogSequenceNumber) (SegmentName, file.File, error) {
	if err := l.init(); err != nil {
		return SegmentName{}, nil, err
	}

	sn := l.getNextSegment(lsn)
	f, err := l.sf.Open(sn)
	if err != nil {
		return SegmentName{}, nil, err
	}

	l.names.PushBack(sn)

	return sn, f, err
}

func (l *List) OpenSegmentBeforeLSN(lsn record.LogSequenceNumber) (SegmentName, file.File, error) {
	if err := l.init(); err != nil {
		return SegmentName{}, nil, err
	}

	sn, _, ok := l.getSegmentBeforeLSN(lsn)
	if !ok {
		return SegmentName{}, nil, ErrNoSegmentFile
	}

	f, err := l.sf.Open(sn)
	if err != nil {
		return SegmentName{}, nil, err
	}

	return sn, f, nil
}

func (l *List) IterateFromLSN(lsn record.LogSequenceNumber) iter.Seq2[file.File, error] {
	return func(yield func(file.File, error) bool) {
		if err := l.init(); err != nil {
			yield(nil, err)
			return
		}

		_, e, ok := l.getSegmentForLSN(lsn)
		if !ok {
			return
		}

		for {
			if e == nil {
				return
			}

			sn := e.Value.(SegmentName)
			f, err := l.sf.Open(sn)
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(f, nil) {
				return
			}

			e = e.Next()
		}
	}
}

func (l *List) OpenSegmentForLSN(lsn record.LogSequenceNumber) (SegmentName, file.File, error) {
	if err := l.init(); err != nil {
		return SegmentName{}, nil, err
	}

	sn, _, ok := l.getSegmentForLSN(lsn)
	if !ok {
		return SegmentName{}, nil, ErrNoSegmentFile
	}

	f, err := l.sf.Open(sn)
	if err != nil {
		return SegmentName{}, nil, err
	}

	return sn, f, err
}

func (l *List) RemoveSegmentBeforeLSN(lsn record.LogSequenceNumber) error {
	if err := l.init(); err != nil {
		return err
	}

	sn, e, ok := l.getSegmentBeforeLSN(lsn)
	if !ok {
		return ErrNoSegmentFile
	}

	if err := l.sf.Remove(sn); err != nil {
		return err
	}

	l.names.Remove(e)
	return nil
}

func (l *List) Reverse() iter.Seq2[file.File, error] {
	return func(yield func(file.File, error) bool) {
		if err := l.init(); err != nil {
			yield(nil, err)
			return
		}

		for sn := range l.reverse() {
			f, err := l.sf.Open(sn)
			if err != nil {
				yield(nil, err)
				return
			}

			if !yield(f, nil) {
				return
			}
		}
	}
}

func (*List) filesToNames(files []fs.DirEntry) *list.List {
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

func (l *List) getFiles() ([]fs.DirEntry, error) {
	return l.fsys.ReadDir(l.dir)
}

func (l *List) getLatestSegment() (SegmentName, bool) {
	e := l.names.Back()
	if e == nil {
		return SegmentName{}, false
	}

	return e.Value.(SegmentName), true
}

func (l *List) getNextSegment(lsn record.LogSequenceNumber) SegmentName {
	sn, ok := l.getLatestSegment()
	if !ok {
		return SegmentName{}
	}

	return sn.Next(lsn)
}

func (l *List) getSegmentBeforeLSN(lsn record.LogSequenceNumber) (SegmentName, *list.Element, bool) {
	for n, e := range l.reverse() {
		if lsn.Value() >= n.lsn.Value() {
			e = e.Prev()
			if e == nil {
				return SegmentName{}, nil, false
			}
			n = e.Value.(SegmentName)
			return n, e, true
		}
	}

	return SegmentName{}, nil, false
}

func (l *List) getSegmentForLSN(lsn record.LogSequenceNumber) (SegmentName, *list.Element, bool) {
	for n, e := range l.reverse() {
		if lsn.Value() >= n.lsn.Value() {
			return n, e, true
		}
	}

	return SegmentName{}, nil, false
}

func (l *List) iterate() iter.Seq2[SegmentName, *list.Element] {
	return func(yield func(SegmentName, *list.Element) bool) {
		if l.names == nil {
			return
		}

		e := l.names.Front()
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

func (l *List) reverse() iter.Seq2[SegmentName, *list.Element] {
	return func(yield func(SegmentName, *list.Element) bool) {
		if l.names == nil {
			return
		}

		e := l.names.Back()
		for {
			if e == nil {
				return
			}

			if !yield(e.Value.(SegmentName), e) {
				return
			}

			e = e.Prev()
		}
	}
}
