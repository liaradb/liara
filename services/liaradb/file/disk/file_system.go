package disk

import (
	"io/fs"
	"os"

	"github.com/liaradb/liaradb/file"
)

type FileSystem struct {
}

func New() *file.Cache {
	return file.NewCache(&FileSystem{})
}

func (fs *FileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (fs *FileSystem) MkDirAll(name string) error {
	return os.MkdirAll(name, 0750)
}

func (fs *FileSystem) OpenFile(name string) (file.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
}

func (fs *FileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (fs *FileSystem) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
