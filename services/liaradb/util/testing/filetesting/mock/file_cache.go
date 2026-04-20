package mock

import (
	"testing/fstest"

	"github.com/liaradb/liaradb/file"
)

func NewFileSystem(fsys fstest.MapFS) file.FileSystem {
	return file.NewFileSystem(newFileSystem(fsys))
}
