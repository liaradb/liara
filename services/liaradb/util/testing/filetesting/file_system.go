package filetesting

import (
	"io/fs"
	"path"
	"sync"
	"testing/fstest"
	"time"

	"github.com/liaradb/liaradb/file"
)

type FileSystem struct {
	fstest.MapFS
	dirs  map[string]map[string]*File
	mux   sync.Mutex
	lock  chan struct{}
	delay time.Duration
}

func New(fsys fstest.MapFS) *file.Cache {
	return file.NewCache(newFileSystem(fsys))
}

func newFileSystem(fsys fstest.MapFS) *FileSystem {
	lock := make(chan struct{})
	close(lock)

	return &FileSystem{
		MapFS: fsys,
		lock:  lock,
	}
}

func NewCacheDelay(fsys fstest.MapFS, delay time.Duration) *file.Cache {
	return file.NewCache(newFileSystemDelay(fsys, delay))
}

func newFileSystemDelay(fsys fstest.MapFS, delay time.Duration) *FileSystem {
	lock := make(chan struct{})
	close(lock)

	return &FileSystem{
		MapFS: fsys,
		lock:  lock,
		delay: delay,
	}
}

func (mfs *FileSystem) Lock()   { mfs.lock = make(chan struct{}) }
func (mfs *FileSystem) UnLock() { close(mfs.lock) }

func (mfs *FileSystem) MkDirAll(name string) error {
	mfs.mux.Lock()
	defer mfs.mux.Unlock()

	if mfs.MapFS == nil {
		mfs.MapFS = make(fstest.MapFS)
	}

	_, err := mfs.MapFS.Stat(name)
	if err == nil {
		return nil
	}

	mfs.MapFS[name] = &fstest.MapFile{
		Mode: fs.ModeDir,
	}

	return nil
}

func (mfs *FileSystem) OpenFile(name string) (file.File, error) {
	mfs.mux.Lock()
	defer mfs.mux.Unlock()

	<-mfs.lock

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
		m = NewMockFile(name, mfs.delay)
		f, ok := mfs.MapFS[dir]
		if ok {
			m.Data = f.Data
			m.ModTime = f.ModTime
		}
		d[base] = m
	} else {
		// Clone
		d[base] = m.clone()
	}

	mfs.MapFS[name] = &m.MapFile

	m.Open()

	if mfs.delay != 0 {
		time.Sleep(mfs.delay)
	}

	return m, nil
}

func (mfs *FileSystem) Remove(name string) error {
	mfs.mux.Lock()
	defer mfs.mux.Unlock()

	d, ok := mfs.dirs[path.Dir(name)]
	if ok {
		delete(d, path.Base(name))
	}

	delete(mfs.MapFS, name)
	return nil
}
