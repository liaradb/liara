package filecache

import (
	"io"
	"io/fs"
)

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
