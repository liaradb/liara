package mock

import (
	"errors"
	"io"
	"os"
	"slices"
	"time"

	"github.com/liaradb/liaradb/file"
)

type File struct {
	name     string
	data     []byte
	position int64
	modTime  time.Time
	isOpen   bool
}

var _ file.File = (*File)(nil)
var ErrClosed = errors.New("file already closed")

func NewMockFile(name string) *File {
	return &File{
		name:    name,
		modTime: time.Now(),
	}
}

func (f *File) Open() {
	f.isOpen = true
}

func (f *File) IsOpen() bool {
	return f.isOpen
}

func (f *File) Close() error {
	f.isOpen = false
	return nil
}

func (f *File) Stat() (os.FileInfo, error) {
	// TODO: Is this correct?
	if !f.isOpen {
		return nil, ErrClosed
	}

	return &mockFileInfo{
		name:    f.name,
		size:    int64(len(f.data)),
		modTime: f.modTime,
	}, nil
}

// TODO: Test this
func (f *File) Read(b []byte) (n int, err error) {
	return f.ReadAt(b, f.position)
}

func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	if !f.isOpen {
		return 0, ErrClosed
	}

	if off < 0 {
		return 0, &os.PathError{
			Op:   "readat",
			Path: f.name,
			Err:  errors.New("negative offset")}
	}

	if f.endOfFile(b, off) {
		err = io.EOF
	}

	n = copy(b, f.data[off:])
	f.position = off + int64(len(b))

	return
}

// TODO: Test this
func (f *File) Write(b []byte) (n int, err error) {
	return f.WriteAt(b, f.position)
}

func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	if !f.isOpen {
		return 0, ErrClosed
	}

	if off < 0 {
		return 0, &os.PathError{
			Op:   "writeat",
			Path: f.name,
			Err:  errors.New("negative offset")}
	}

	f.modTime = time.Now()
	f.adjustSize(b, off)
	f.position = off + int64(len(b))

	return copy(f.data[off:int(off)+len(b)], b), nil
}

func (f *File) endOfFile(b []byte, off int64) bool {
	l := int(off) + len(b)
	g := l - len(f.data)
	return g > 0
}

func (f *File) adjustSize(b []byte, off int64) bool {
	l := int(off) + len(b)
	g := l - len(f.data)
	if g > 0 {
		f.data = slices.Grow(f.data, g)
		f.data = f.data[0:l]
		return true
	}
	return false
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if !f.isOpen {
		return 0, ErrClosed
	}

	f.position = f.seekPosition(offset, whence)
	return f.position, nil
}

func (f *File) seekPosition(offset int64, whence int) int64 {
	switch whence {
	case io.SeekStart:
		return offset
	case io.SeekCurrent:
		return f.position + offset
	case io.SeekEnd:
		return int64(len(f.data)) - offset
	default:
		return offset
	}
}
