package file

import (
	"io/fs"
)

type FileSystem interface {
	MkDirAll(name string) error
	OpenFile(name string) (File, error)
	Remove(name string) error
	ReadDir(name string) ([]fs.DirEntry, error)
	Stat(name string) (fs.FileInfo, error)
}
