package segment

import (
	"container/list"
	"io/fs"
	"iter"

	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
)

type List struct {
	dir         string
	fsys        filecache.FileSystem
	names       *list.List
	f           *File
	pageSize    int64
	segmentSize action.PageID
}

func NewList(
	fsys filecache.FileSystem,
	dir string,
	pageSize int64,
	segmentSize action.PageID,
) *List {
	return &List{
		dir:         dir,
		fsys:        fsys,
		pageSize:    pageSize,
		segmentSize: segmentSize,
	}
}

func (l *List) file() (*File, bool) {
	if l.f == nil {
		return nil, false
	}

	return l.f, true
}

func (l *List) path(sn SegmentName) string {
	return sn.path(l.dir)
}

func (l *List) Names() []SegmentName {
	names := make([]SegmentName, 0, l.names.Len())
	for n := range l.iterate() {
		names = append(names, n)
	}
	return names
}

func (l *List) Close() error {
	f, ok := l.file()
	if !ok || !f.IsOpen() {
		return nil
	}

	if err := f.close(); err != nil {
		return err
	}

	l.f = nil
	return nil
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

func (l *List) OpenLatestSegment() (*File, error) {
	if err := l.init(); err != nil {
		return nil, err
	}

	sn, ok := l.getLatestSegment()
	f, err := l.open(sn)
	if err != nil {
		return nil, err
	}

	if !ok {
		l.names.PushBack(sn)
	}

	return f, err
}

func (l *List) OpenNextSegment(lsn record.LogSequenceNumber) (*File, error) {
	if err := l.init(); err != nil {
		return nil, err
	}

	sn := l.getNextSegment(lsn)
	f, err := l.open(sn)
	if err != nil {
		return nil, err
	}

	l.names.PushBack(sn)

	return f, err
}

func (l *List) OpenSegmentBeforeLSN(lsn record.LogSequenceNumber) (*File, error) {
	if err := l.init(); err != nil {
		return nil, err
	}

	sn, _, ok := l.getSegmentBeforeLSN(lsn)
	if !ok {
		return nil, ErrNoSegmentFile
	}

	f, err := l.open(sn)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (l *List) IterateFromLSN(lsn record.LogSequenceNumber) iter.Seq2[*File, error] {
	return func(yield func(*File, error) bool) {
		if err := l.init(); err != nil {
			yield(nil, err)
			return
		}

		_, e, ok := l.getSegmentForLSN(lsn)
		if !ok {
			return
		}

		for sn := range iterate[SegmentName](e) {
			if f, err := l.open(sn); !yield(f, err) {
				return
			}
		}
	}
}

func (l *List) OpenSegmentForLSN(lsn record.LogSequenceNumber) (*File, error) {
	if err := l.init(); err != nil {
		return nil, err
	}

	sn, _, ok := l.getSegmentForLSN(lsn)
	if !ok {
		return nil, ErrNoSegmentFile
	}

	f, err := l.open(sn)
	if err != nil {
		return nil, err
	}

	return f, err
}

func (l *List) open(sn SegmentName) (*File, error) {
	if f, ok := l.isCurrentAndOpen(sn); ok {
		return f, nil
	}

	if err := l.Close(); err != nil {
		return nil, err
	}

	return l.openFile(sn)
}

func (l *List) openFile(sn SegmentName) (*File, error) {
	f, err := l.fsys.OpenFile(l.path(sn))
	if err != nil {
		return nil, err
	}

	file := newFile(f, sn, l.pageSize, l.segmentSize)
	l.f = file
	return file, nil
}

func (l *List) isCurrentAndOpen(sn SegmentName) (*File, bool) {
	if f, ok := l.file(); ok && f.isCurrentAndOpen(sn) {
		return f, true

	}

	return nil, false
}

func (l *List) RemoveSegmentBeforeLSN(lsn record.LogSequenceNumber) error {
	if err := l.init(); err != nil {
		return err
	}

	sn, e, ok := l.getSegmentBeforeLSN(lsn)
	if !ok {
		return ErrNoSegmentFile
	}

	if err := l.remove(sn); err != nil {
		return err
	}

	l.names.Remove(e)
	return nil
}

func (l *List) remove(sn SegmentName) error {
	if f, ok := l.file(); ok {
		if err := f.close(); err != nil {
			return err
		}
	}

	return l.fsys.Remove(l.path(sn))
}

func (l *List) Reverse() iter.Seq2[*File, error] {
	return func(yield func(*File, error) bool) {
		if err := l.init(); err != nil {
			yield(nil, err)
			return
		}

		for sn := range l.reverse() {
			if f, err := l.open(sn); !yield(f, err) {
				return
			}
		}
	}
}

// files are assumed to be sorted
func (*List) filesToNames(files []fs.DirEntry) *list.List {
	names := make([]SegmentName, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, ParseSegmentName(f.Name()))
		}
	}

	// slices.SortFunc(names, func(a, b SegmentName) int {
	// 	return int(a.ID() - b.ID())
	// })

	l := list.New()
	for _, n := range names {
		l.PushBack(n)
	}

	return l
}

func (l *List) getFiles() ([]fs.DirEntry, error) {
	if err := l.fsys.MkDirAll(l.dir); err != nil {
		return nil, err
	}

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

func (l *List) iterate() iter.Seq[SegmentName] {
	if l.names == nil {
		return iterate[SegmentName](nil)
	}

	return iterate[SegmentName](l.names.Front())
}

func iterate[T any](e *list.Element) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			if e == nil {
				return
			}

			if t := e.Value.(T); !yield(t) {
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
