package file

import (
	"errors"
	"io/fs"
	"sync"
)

type FileSsytemCache struct {
	files map[string]*file
	mux   sync.RWMutex
	fsys  FileSystem
}

func (fsc *FileSsytemCache) MkDirAll(name string) error {
	return fsc.fsys.MkDirAll(name)
}

func (fsc *FileSsytemCache) OpenFile(name string) (File, error) {
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
		fsc.files = map[string]*file{}
	}
	df = &file{
		File: f,
		name: name,
		fsys: fsc}
	fsc.files[name] = df
	return df, nil
}

func (fsc *FileSsytemCache) ReadDir(name string) ([]fs.DirEntry, error) {
	return fsc.fsys.ReadDir(name)
}

func (fsc *FileSsytemCache) Remove(name string) error {
	if err := fsc.CloseFile(name); err != nil {
		return err
	}

	return fsc.fsys.Remove(name)
}

func (fsc *FileSsytemCache) Count() int {
	fsc.mux.RLock()
	defer fsc.mux.RUnlock()

	return len(fsc.files)
}

func (fsc *FileSsytemCache) Close() error {
	errs := make([]error, 0, len(fsc.files))
	for n := range fsc.files {
		errs = append(errs, fsc.CloseFile(n))
	}
	return errors.Join(errs...)
}

func (fsc *FileSsytemCache) CloseFile(name string) error {
	fsc.mux.Lock()
	defer fsc.mux.Unlock()

	f, ok := fsc.files[name]
	if !ok {
		return nil
	}

	if err := f.Close(); err != nil {
		return err
	}

	delete(fsc.files, name)
	return nil
}
