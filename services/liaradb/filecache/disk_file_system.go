package filecache

import (
	"io/fs"
	"os"
)

type DiskFileSystem struct {
}

func (fs *DiskFileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (fs *DiskFileSystem) MkDirAll(name string) error {
	return os.MkdirAll(name, 0750)
}

func (fs *DiskFileSystem) OpenFile(name string) (File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
}

func (fs *DiskFileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (fs *DiskFileSystem) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
