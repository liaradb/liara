package filetesting

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"slices"
	"testing/fstest"
	"time"

	"github.com/liaradb/liaradb/filecache"
)

type File struct {
	*fileInner
	position int64
}

type fileInner struct {
	fstest.MapFile
	name       string
	isOpen     bool
	readCount  int
	writeCount int
	delay      time.Duration
	fsys       *FileSystem
}

var _ filecache.File = (*File)(nil)

func NewMockFile(name string, delay time.Duration, mod time.Time) *File {
	return NewMockFileWithFileSystem(name, delay, nil, mod)
}

func NewMockFileWithFileSystem(name string, delay time.Duration, fsys *FileSystem, mod time.Time) *File {
	return &File{
		fileInner: &fileInner{
			name: name,
			MapFile: fstest.MapFile{
				ModTime: mod,
			},
			delay: delay,
			fsys:  fsys,
		},
	}
}

func (f *File) clone() *File {
	return &File{
		fileInner: f.fileInner,
	}
}

func (f *File) ReadCount() int  { return f.readCount }
func (f *File) WriteCount() int { return f.writeCount }

func (f *File) Open() {
	f.isOpen = true
	f.position = 0
}

func (f *File) IsOpen() bool {
	return f.isOpen
}

func (f *File) Close() error {
	f.isOpen = false
	return nil
}

func (f *File) Stat() (os.FileInfo, error) {
	if !f.isOpen {
		return nil, fs.ErrClosed
	}

	return &mockFileInfo{
		name:    f.name,
		size:    int64(len(f.Data)),
		modTime: f.ModTime,
	}, nil
}

func (f *File) Read(b []byte) (n int, err error) {
	return f.ReadAt(b, f.position)
}

func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	if !f.isOpen {
		return 0, fs.ErrClosed
	}

	if off < 0 {
		return 0, &os.PathError{
			Op:   "readat",
			Path: f.name,
			Err:  errors.New("negative offset")}
	}

	if f.isFail() {
		return 0, io.ErrUnexpectedEOF
	}

	if f.endOfFile(b, off) {
		err = io.EOF
	}

	if f.delay != 0 {
		time.Sleep(f.delay)
	}

	if off <= int64(len(f.Data)) {
		n = copy(b, f.Data[off:])
	}

	f.position = off + int64(len(b))
	f.readCount++

	return
}

func (f *File) Write(b []byte) (n int, err error) {
	return f.WriteAt(b, f.position)
}

func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	if !f.isOpen {
		return 0, fs.ErrClosed
	}

	if off < 0 {
		return 0, &os.PathError{
			Op:   "writeat",
			Path: f.name,
			Err:  errors.New("negative offset")}
	}

	if f.isFail() {
		return 0, io.ErrUnexpectedEOF
	}

	if f.delay != 0 {
		time.Sleep(f.delay)
	}

	f.MapFile.ModTime = time.Now()
	f.adjustSize(b, off)
	f.position = off + int64(len(b))
	f.writeCount++

	return copy(f.Data[off:int(off)+len(b)], b), nil
}

func (f *File) isFail() bool {
	if f.fsys == nil {
		return false
	}

	return f.fsys.fail
}

func (f *File) endOfFile(b []byte, off int64) bool {
	l := int(off) + len(b)
	g := l - len(f.Data)
	return g > 0
}

func (f *File) adjustSize(b []byte, off int64) bool {
	l := int(off) + len(b)
	g := l - len(f.Data)
	if g > 0 {
		f.Data = slices.Grow(f.Data, g)
		f.Data = f.Data[0:l]
		return true
	}
	return false
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if !f.isOpen {
		return 0, fs.ErrClosed
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
		return int64(len(f.Data)) - offset
	default:
		return offset
	}
}
