package segment

import (
	"io"
	"io/fs"

	"github.com/liaradb/liaradb/filecache"
)

type File struct {
	file        filecache.File
	sn          SegmentName
	pageSize    int64
	segmentSize PageID
	pageID      PageID
	size        int64
}

func newFile(
	file filecache.File,
	sn SegmentName,
	pageSize int64,
	segmentSize PageID,
	size int64,
) *File {
	return &File{
		file:        file,
		sn:          sn,
		pageSize:    pageSize,
		segmentSize: segmentSize,
		size:        size,
	}
}

func (f *File) SegmentName() SegmentName { return f.sn }
func (f *File) Size() int64              { return f.size }

func (f *File) isCurrent(sn SegmentName) bool {
	return f.sn == sn
}

func (f *File) isCurrentAndOpen(sn SegmentName) bool {
	return f.isCurrent(sn) && f.IsOpen()
}

func (f *File) IsOpen() bool {
	return f.file != nil
}

func (f *File) Close() error {
	if !f.IsOpen() {
		return nil
	}

	if err := f.file.Close(); err != nil {
		return err
	}

	f.file = nil
	return nil
}

func (f *File) refreshSize() (int64, error) {
	stat, err := f.file.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

func (f *File) Stat() (fs.FileInfo, error) {
	return f.file.Stat()
}

func (f *File) SeekTail() (int64, error) {
	size, err := f.refreshSize()
	if err != nil {
		return 0, err
	}

	f.pageID = newActivePageIDFromSize(size, f.pageSize)
	return size, nil
}

func (f *File) ReadAt(data []byte, off int64) (int, error) {
	return f.file.ReadAt(data, off)
}

func (f *File) WriteAt(data []byte, off int64) (int, error) {
	return f.file.WriteAt(data, off)
}

func (f *File) Read(data []byte) (int, error) {
	wr := io.NewSectionReader(f, f.position(), f.pageSize)
	return wr.Read(data)
}

func (f *File) Write(data []byte) error {
	off := f.position()
	wr := io.NewOffsetWriter(f, off)
	n, err := wr.Write(data)
	f.size = max(f.size, off+int64(n))
	return err
}

func (f *File) IsEmpty() bool {
	return f.size == 0
}

func (f *File) NextPage() bool {
	if f.pageID+1 >= f.segmentSize {
		return false
	}

	f.pageID++
	f.size += f.pageSize
	return true
}

func (f *File) NextPageUntilSize(size int64) bool {
	segmentSize := newActivePageIDFromSize(size, f.pageSize)
	if f.pageID >= segmentSize {
		return false
	}

	f.pageID++
	return true
}

func (f *File) PrevPageUntilStart() bool {
	if f.pageID <= 0 {
		return false
	}

	f.pageID--
	return true
}

func (f *File) position() int64 {
	return f.pageID.Position(f.pageSize)
}

func (f *File) reset() {
	f.pageID = 0
}
