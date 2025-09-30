package disk

import "github.com/liaradb/liaradb/file"

type File struct {
	file.File
	name string
	fsys *FileSystem
}

// TODO: Test this
func (f *File) Close() error {
	return f.fsys.CloseFile(f.name)
}
