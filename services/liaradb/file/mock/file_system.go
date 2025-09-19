package mock

import (
	"path"
	"testing/fstest"

	"github.com/liaradb/liaradb/file"
)

type FileSystem struct {
	fstest.MapFS
	dirs map[string]map[string]*mockFile
}

func NewFileSystem(fs fstest.MapFS) *FileSystem {
	return &FileSystem{
		MapFS: fs,
	}
}

func (mfs *FileSystem) OpenFile(name string) (file.File, error) {
	if mfs.MapFS == nil {
		mfs.MapFS = make(fstest.MapFS)
	}
	if mfs.dirs == nil {
		mfs.dirs = make(map[string]map[string]*mockFile)
	}

	dir := path.Dir(name)
	d, ok := mfs.dirs[dir]
	if !ok {
		d = make(map[string]*mockFile)
		mfs.dirs[dir] = d
	}

	m, ok := d[name]
	if !ok {
		m = NewMockFile(name)
		f, ok := mfs.MapFS[dir]
		if ok {
			m.data = f.Data
			m.modTime = f.ModTime
		}
		d[name] = m
	}

	mfs.MapFS[name] = &fstest.MapFile{
		Data:    m.data,
		ModTime: m.modTime,
	}

	return m, nil
}
