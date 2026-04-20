package filecache

import "github.com/liaradb/liaradb/file"

type File struct {
	file.File
	name   string
	fsys   *Cache
	closed bool
}

func (f *File) Close() error {
	return f.fsys.CloseFile(f.name)
}

func (f *File) closeFile() error {
	err := f.File.Close()
	if err != nil {
		return err
	}

	f.closed = true
	return nil
}

func (f *File) IsOpen() bool {
	return !f.closed
}
