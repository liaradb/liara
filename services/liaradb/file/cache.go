package file

import (
	"errors"
	"io/fs"
	"sync"
)

type Cache struct {
	files map[string]*CacheFile
	mux   sync.RWMutex
	fsys  FileSystem
}

func NewCache(
	fsys FileSystem,
) *Cache {
	return &Cache{
		fsys: fsys,
	}
}

func (fsc *Cache) MkDirAll(name string) error {
	return fsc.fsys.MkDirAll(name)
}

func (fsc *Cache) OpenFile(name string) (File, error) {
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
		fsc.files = map[string]*CacheFile{}
	}
	df = &CacheFile{
		File: f,
		name: name,
		fsys: fsc}
	fsc.files[name] = df
	return df, nil
}

func (fsc *Cache) ReadDir(name string) ([]fs.DirEntry, error) {
	return fsc.fsys.ReadDir(name)
}

func (fsc *Cache) Stat(name string) (fs.FileInfo, error) {
	return fsc.fsys.Stat(name)
}

func (fsc *Cache) Remove(name string) error {
	if err := fsc.CloseFile(name); err != nil {
		return err
	}

	return fsc.fsys.Remove(name)
}

func (fsc *Cache) Count() int {
	fsc.mux.RLock()
	defer fsc.mux.RUnlock()

	return len(fsc.files)
}

func (fsc *Cache) Close() error {
	errs := make([]error, 0, len(fsc.files))
	for n := range fsc.files {
		errs = append(errs, fsc.CloseFile(n))
	}
	return errors.Join(errs...)
}

func (fsc *Cache) CloseFile(name string) error {
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
