package storage

import (
	"errors"
	"io"
	"io/fs"
	"os"
)

type FS interface {
	Open(name string) (file, error)
}

type file interface {
	fs.File
	io.ReaderAt
	io.WriterAt
	io.Seeker
}

type fileSystem struct {
	files map[string]file
}

func (fs *fileSystem) Close() error {
	errs := make([]error, 0, len(fs.files))
	for _, f := range fs.files {
		errs = append(errs, f.Close())
	}
	return errors.Join(errs...)
}

func (fs *fileSystem) Open(name string) (file, error) {
	f, ok := fs.files[name]
	if ok {
		return f, nil
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	if fs.files == nil {
		fs.files = map[string]file{}
	}
	fs.files[name] = f
	return f, nil
}

func (fs *fileSystem) CloseFile(name string) error {
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
