package segment

import (
	"io"
	"io/fs"

	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/action"
)

type File struct {
	file        filecache.File
	sn          SegmentName
	pageSize    int64
	segmentSize action.PageID
	pageID      action.PageID
}

func newFile(
	file filecache.File,
	sn SegmentName,
	pageSize int64,
	segmentSize action.PageID,
) *File {
	return &File{
		file:        file,
		sn:          sn,
		pageSize:    pageSize,
		segmentSize: segmentSize,
	}
}

func (f *File) SegmentName() SegmentName { return f.sn }

func (f *File) isCurrent(sn SegmentName) bool {
	return f.sn == sn
}

func (f *File) isCurrentAndOpen(sn SegmentName) bool {
	return f.isCurrent(sn) && f.IsOpen()
}

func (f *File) IsOpen() bool {
	return f.file != nil
}

func (f *File) close() error {
	if !f.IsOpen() {
		return nil
	}

	if err := f.file.Close(); err != nil {
		return err
	}

	f.file = nil
	return nil
}

func (f *File) Size() (int64, error) {
	stat, err := f.file.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

func (f *File) Stat() (fs.FileInfo, error) {
	return f.file.Stat()
}

func (f *File) SeekTail() error {
	size, err := f.Size()
	if err != nil {
		return err
	}

	f.pageID = action.NewActivePageIDFromSize(size, f.pageSize)
	return nil
}

func (f *File) ReadAt(data []byte, off int64) (int, error) {
	return f.file.ReadAt(data, off)
}

func (f *File) WriteAt(data []byte, off int64) (int, error) {
	return f.file.WriteAt(data, off)
}

func (ps *File) Write(data []byte) error {
	wr := io.NewOffsetWriter(ps, ps.position())
	_, err := wr.Write(data)
	return err
}

func (f *File) IncrementPageID() bool {
	if f.pageID+1 >= f.segmentSize {
		return false
	}

	f.pageID++
	return true
}

func (f *File) position() int64 {
	return f.pageID.Position(f.pageSize)
}
