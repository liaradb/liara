package mock

import (
	"path"
	"testing/fstest"

	"github.com/liaradb/liaradb/file"
)

type FileSystem struct {
	fstest.MapFS
	dirs map[string]map[string]*File
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
		mfs.dirs = make(map[string]map[string]*File)
	}

	dir := path.Dir(name)
	d, ok := mfs.dirs[dir]
	if !ok {
		d = make(map[string]*File)
		mfs.dirs[dir] = d
	}

	base := path.Base(name)
	m, ok := d[base]
	if !ok {
		m = NewMockFile(name)
		m.Open()
		f, ok := mfs.MapFS[dir]
		if ok {
			m.data = f.Data
			m.modTime = f.ModTime
		}
		d[base] = m
	}

	mfs.MapFS[name] = &fstest.MapFile{
		Data:    m.data,
		ModTime: m.modTime,
	}

	return m, nil
}

func (mfs *FileSystem) Remove(name string) error {
	d, ok := mfs.dirs[path.Dir(name)]
	if ok {
		delete(d, path.Base(name))
	}

	delete(mfs.MapFS, name)
	return nil
}
