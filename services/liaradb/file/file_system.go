package file

import (
	"io"
	"io/fs"
)

type FS interface {
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
