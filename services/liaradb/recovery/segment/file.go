package segment

import (
	"io/fs"

	"github.com/liaradb/liaradb/filecache"
)

type File struct {
	file filecache.File
	sn   SegmentName
}

func newFile(
	file filecache.File,
	sn SegmentName,
) *File {
	return &File{
		file: file,
		sn:   sn,
	}
}

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

func (f *File) ReadAt(data []byte, off int64) (int, error) {
	return f.file.ReadAt(data, off)
}

func (f *File) WriteAt(data []byte, off int64) (int, error) {
	return f.file.WriteAt(data, off)
}
