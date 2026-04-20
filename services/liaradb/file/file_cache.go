package file

import (
	"errors"
	"io/fs"
	"sync"
)

type FileCache struct {
	files map[string]*FileCacheFile
	mux   sync.RWMutex
	fsys  FileSystem
}

func NewFileSystem(
	fsys FileSystem,
) *FileCache {
	return &FileCache{
		fsys: fsys,
	}
}

func (fsc *FileCache) MkDirAll(name string) error {
	return fsc.fsys.MkDirAll(name)
}

func (fsc *FileCache) OpenFile(name string) (File, error) {
	fsc.mux.Lock()
	defer fsc.mux.Unlock()

	df, ok := fsc.files[name]
	if ok {
		return df, nil
	}

	f, err := fsc.fsys.OpenFile(name)
	if err != nil {
		return nil, err
	}

	if fsc.files == nil {
		fsc.files = map[string]*FileCacheFile{}
	}
	df = &FileCacheFile{
		File: f,
		name: name,
		fsys: fsc}
	fsc.files[name] = df
	return df, nil
}

func (fsc *FileCache) ReadDir(name string) ([]fs.DirEntry, error) {
	return fsc.fsys.ReadDir(name)
}

func (fsc *FileCache) Stat(name string) (fs.FileInfo, error) {
	return fsc.fsys.Stat(name)
}

func (fsc *FileCache) Remove(name string) error {
	if err := fsc.CloseFile(name); err != nil {
		return err
	}

	return fsc.fsys.Remove(name)
}

func (fsc *FileCache) Count() int {
	fsc.mux.RLock()
	defer fsc.mux.RUnlock()

	return len(fsc.files)
}

func (fsc *FileCache) Close() error {
	errs := make([]error, 0, len(fsc.files))
	for n := range fsc.files {
		errs = append(errs, fsc.CloseFile(n))
	}
	return errors.Join(errs...)
}

func (fsc *FileCache) CloseFile(name string) error {
	fsc.mux.Lock()
	defer fsc.mux.Unlock()

	f, ok := fsc.files[name]
	if !ok {
		return nil
	}

	if err := f.closeFile(); err != nil {
		return err
	}

	delete(fsc.files, name)
	return nil
}
