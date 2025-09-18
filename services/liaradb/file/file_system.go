package file

import (
	"io"
	"io/fs"
)

type FileSystem interface {
	ReadDir(name string) ([]fs.DirEntry, error)
	Open(name string) (File, error)
}

type File interface {
	fs.File
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
	io.Seeker
}
