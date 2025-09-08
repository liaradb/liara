package file

import (
	"errors"
	"io"
	"io/fs"
	"os"
)

type File interface {
	fs.File
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
	io.Seeker
}

type FileSystem struct {
	files map[string]File
}

func (fs *FileSystem) Close() error {
	errs := make([]error, 0, len(fs.files))
	for n := range fs.files {
		errs = append(errs, fs.CloseFile(n))
	}
	return errors.Join(errs...)
}

func (fs *FileSystem) Open(name string) (File, error) {
	f, ok := fs.files[name]
	if ok {
		return f, nil
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	if fs.files == nil {
		fs.files = map[string]File{}
	}
	fs.files[name] = f
	return f, nil
}

func (fs *FileSystem) CloseFile(name string) error {
	f, ok := fs.files[name]
	if !ok {
		return nil
	}

	if err := f.Close(); err != nil {
		return err
	}

	delete(fs.files, name)
	return nil
}

func (fs *FileSystem) Count() int {
	return len(fs.files)
}
