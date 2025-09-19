package mock

import (
	"github.com/liaradb/liaradb/file"
)

type MockFileSource struct {
	dirs map[string]map[string]*mockFile
}

func (mfs *MockFileSource) OpenFile(dir string, name string) (file.File, error) {
	if mfs.dirs == nil {
		mfs.dirs = make(map[string]map[string]*mockFile)
	}

	d, ok := mfs.dirs[dir]
	if !ok {
		d = make(map[string]*mockFile)
		mfs.dirs[dir] = d
	}

	m, ok := d[name]
	if !ok {
		m = NewMockFile(name)
		d[name] = m
	}

	return m, nil
}

func (*MockFileSource) InitDirectory(dir string) error {
	return nil
}
