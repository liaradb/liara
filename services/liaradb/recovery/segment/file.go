package segment

import (
	"path"

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
	return f.isCurrent(sn) && f.isOpen()
}

func (f *File) isOpen() bool {
	return f.file != nil
}

func (f *File) path(dir string) string {
	return path.Join(dir, f.sn.String())
}

func (f *File) close() error {
	if !f.isOpen() {
		return nil
	}

	if err := f.file.Close(); err != nil {
		return err
	}

	f.file = nil
	return nil
}
