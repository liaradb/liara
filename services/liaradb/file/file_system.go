package file

import (
	"io"
	"io/fs"
)

type FileSystem interface {
	MkDirAll(name string) error
	OpenFile(name string) (File, error)
	Remove(name string) error
	ReadDir(name string) ([]fs.DirEntry, error)
	Stat(name string) (fs.FileInfo, error)
}

type File interface {
	fs.File
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
	io.Seeker
	Close
	Stat
}

type Close interface {
	Close() error
}

type Stat interface {
	Stat() (fs.FileInfo, error)
}
