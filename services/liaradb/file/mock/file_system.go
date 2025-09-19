package mock

import (
	"testing/fstest"

	"github.com/liaradb/liaradb/file"
)

type FileSystem struct {
	fstest.MapFS
}

func NewFileSystem(fs fstest.MapFS) *FileSystem {
	return &FileSystem{
		MapFS: fs,
	}
}

func (mf *FileSystem) OpenFile(name string) (file.File, error) {
	return nil, nil
}
