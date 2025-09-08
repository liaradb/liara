package mock

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"slices"
	"time"
)

type mockFile struct {
	name     string
	data     []byte
	position int64
	modTime  time.Time
}

type file interface {
	fs.File
	io.ReaderAt
	io.WriterAt
}

var _ file = (*mockFile)(nil)

func NewMockFile(name string) *mockFile {
	return &mockFile{
		name:    name,
		modTime: time.Now(),
	}
}

func (m *mockFile) Close() error {
	return nil
}

func (m *mockFile) Stat() (os.FileInfo, error) {
	return &mockFileInfo{
		name:    m.name,
		size:    int64(len(m.data)),
		modTime: m.modTime,
	}, nil
}

// TODO: Test this
func (m *mockFile) Read(b []byte) (n int, err error) {
	return m.WriteAt(b, m.position)
}

func (m *mockFile) ReadAt(b []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, &os.PathError{
			Op:   "readat",
			Path: m.name,
			Err:  errors.New("negative offset")}
	}

	m.adjustSize(b, off)

	n = copy(b, m.data[off:])
	m.position = off + int64(len(b))

	return n, nil
}

// TODO: Test this
func (m *mockFile) Write(b []byte) (n int, err error) {
	return m.WriteAt(b, m.position)
}

func (m *mockFile) WriteAt(b []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, &os.PathError{
			Op:   "writeat",
			Path: m.name,
			Err:  errors.New("negative offset")}
	}

	m.modTime = time.Now()
	m.adjustSize(b, off)
	m.position = off + int64(len(b))

	return copy(m.data[off:int(off)+len(b)], b), nil
}

func (m *mockFile) adjustSize(b []byte, off int64) {
	l := int(off) + len(b)
	g := l - len(m.data)
	if g > 0 {
		m.data = slices.Grow(m.data, g)
		m.data = m.data[0:l]
	}
}
