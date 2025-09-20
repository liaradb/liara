package mock

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"slices"
	"time"

	"github.com/liaradb/liaradb/file"
)

type mockFile struct {
	name     string
	data     []byte
	position int64
	modTime  time.Time
	isOpen   bool
}

var _ file.File = (*mockFile)(nil)

func NewMockFile(name string) *mockFile {
	return &mockFile{
		name:    name,
		modTime: time.Now(),
	}
}

func (m *mockFile) Open() {
	m.isOpen = true
}

func (m *mockFile) Close() error {
	m.isOpen = false
	return nil
}

func (m *mockFile) Stat() (os.FileInfo, error) {
	// TODO: Is this correct?
	if !m.isOpen {
		return nil, fs.ErrClosed
	}

	return &mockFileInfo{
		name:    m.name,
		size:    int64(len(m.data)),
		modTime: m.modTime,
	}, nil
}

// TODO: Test this
func (m *mockFile) Read(b []byte) (n int, err error) {
	return m.ReadAt(b, m.position)
}

func (m *mockFile) ReadAt(b []byte, off int64) (n int, err error) {
	if !m.isOpen {
		return 0, fs.ErrClosed
	}

	if off < 0 {
		return 0, &os.PathError{
			Op:   "readat",
			Path: m.name,
			Err:  errors.New("negative offset")}
	}

	if m.endOfFile(b, off) {
		err = io.EOF
	}

	n = copy(b, m.data[off:])
	m.position = off + int64(len(b))

	return
}

// TODO: Test this
func (m *mockFile) Write(b []byte) (n int, err error) {
	return m.WriteAt(b, m.position)
}

func (m *mockFile) WriteAt(b []byte, off int64) (n int, err error) {
	if !m.isOpen {
		return 0, fs.ErrClosed
	}

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

func (m *mockFile) endOfFile(b []byte, off int64) bool {
	l := int(off) + len(b)
	g := l - len(m.data)
	return g > 0
}

func (m *mockFile) adjustSize(b []byte, off int64) bool {
	l := int(off) + len(b)
	g := l - len(m.data)
	if g > 0 {
		m.data = slices.Grow(m.data, g)
		m.data = m.data[0:l]
		return true
	}
	return false
}

func (m *mockFile) Seek(offxset int64, whence int) (int64, error) {
	if !m.isOpen {
		return 0, fs.ErrClosed
	}

	m.position = m.seekPosition(offxset, whence)
	return m.position, nil
}

func (m *mockFile) seekPosition(offset int64, whence int) int64 {
	switch whence {
	case io.SeekStart:
		return offset
	case io.SeekCurrent:
		return m.position + offset
	case io.SeekEnd:
		return int64(len(m.data)) - offset
	default:
		return offset
	}
}
