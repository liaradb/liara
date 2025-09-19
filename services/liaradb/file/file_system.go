package file

import (
	"io"
	"io/fs"
)

type FileSystem interface {
	Remove(name string) error
	OpenFile(name string) (File, error)
	ReadDir(name string) ([]fs.DirEntry, error)
}

type File interface {
	fs.File
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
	io.Seeker
}
