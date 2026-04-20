package file

import (
	"io/fs"
)

type FileCacheFile struct {
	File
	name   string
	fsys   *FileCache
	closed bool
}

func (f *FileCacheFile) Close() error {
	return f.fsys.CloseFile(f.name)
}

func (f *FileCacheFile) closeFile() error {
	err := f.File.Close()
	if err != nil {
		return err
	}

	f.closed = true
	return nil
}

func (f *FileCacheFile) Stat() (fs.FileInfo, error) {
	return f.File.Stat()
}

func (f *FileCacheFile) IsOpen() bool {
	return !f.closed
}
