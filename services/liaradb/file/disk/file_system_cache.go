package disk

import (
	"errors"
	"io/fs"
	"sync"

	"github.com/liaradb/liaradb/file"
)

type FileSystem struct {
	files map[string]*File
	mux   sync.RWMutex
	fsys  file.FileSystem
}

func NewFileSystem(
	fsys file.FileSystem,
) *FileSystem {
	return &FileSystem{
		fsys: fsys,
	}
}

func (fsc *FileSystem) MkDirAll(name string) error {
	return fsc.fsys.MkDirAll(name)
}

func (fsc *FileSystem) OpenFile(name string) (file.File, error) {
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
		fsc.files = map[string]*File{}
	}
	df = &File{
		File: f,
		name: name,
		fsys: fsc}
	fsc.files[name] = df
	return df, nil
}

func (fsc *FileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return fsc.fsys.ReadDir(name)
}

func (fsc *FileSystem) Stat(name string) (fs.FileInfo, error) {
	return fsc.fsys.Stat(name)
}

func (fsc *FileSystem) Remove(name string) error {
	if err := fsc.CloseFile(name); err != nil {
		return err
	}

	return fsc.fsys.Remove(name)
}

func (fsc *FileSystem) Count() int {
	fsc.mux.RLock()
	defer fsc.mux.RUnlock()

	return len(fsc.files)
}

func (fsc *FileSystem) Close() error {
	errs := make([]error, 0, len(fsc.files))
	for n := range fsc.files {
		errs = append(errs, fsc.CloseFile(n))
	}
	return errors.Join(errs...)
}

func (fsc *FileSystem) CloseFile(name string) error {
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
