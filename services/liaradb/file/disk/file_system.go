package disk

import (
	"errors"
	"io/fs"
	"os"

	"github.com/liaradb/liaradb/file"
)

type FileSystem struct {
	files map[string]*File
}

func (fs *FileSystem) Close() error {
	errs := make([]error, 0, len(fs.files))
	for n := range fs.files {
		errs = append(errs, fs.CloseFile(n))
	}
	return errors.Join(errs...)
}

func (fs *FileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (fs *FileSystem) MkDirAll(name string) error {
	return os.MkdirAll(name, 0750)
}

func (fs *FileSystem) OpenFile(name string) (file.File, error) {
	df, ok := fs.files[name]
	if ok {
		return df, nil
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	if fs.files == nil {
		fs.files = map[string]*File{}
	}
	df = &File{
		File: f,
		name: name,
		fsys: fs}
	fs.files[name] = df
	return df, nil
}

func (fs *FileSystem) CloseFile(name string) error {
	f, ok := fs.files[name]
	if !ok {
		return nil
	}

	if err := f.File.Close(); err != nil {
		return err
	}

	delete(fs.files, name)
	return nil
}

func (fs *FileSystem) Count() int {
	return len(fs.files)
}

func (fs *FileSystem) Remove(name string) error {
	return os.Remove(name)
}
