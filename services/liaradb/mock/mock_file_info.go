package mock

import (
	"io/fs"
	"os"
	"time"
)

type mockFileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

var _ os.FileInfo = (*mockFileInfo)(nil)

func (m *mockFileInfo) IsDir() bool {
	return false
}

func (m *mockFileInfo) ModTime() time.Time {
	return m.modTime
}

func (m *mockFileInfo) Mode() fs.FileMode {
	return fs.ModeAppend
}

func (m *mockFileInfo) Name() string {
	return m.name
}

func (m *mockFileInfo) Size() int64 {
	return m.size
}

func (m *mockFileInfo) Sys() any {
	return nil
}
