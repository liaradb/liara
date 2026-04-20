package filecache

import (
	"errors"
	"io/fs"
	"sync"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/disk"
)

type Cache struct {
	files map[string]*File
	mux   sync.RWMutex
	fsys  file.FileSystem
}

func New() *Cache {
	return NewWithFileSystem(&disk.FileSystem{})
}

func NewWithFileSystem(
	fsys file.FileSystem,
) *Cache {
	return &Cache{
		fsys: fsys,
	}
}

func (c *Cache) FSYS() file.FileSystem { return c.fsys }

func (c *Cache) MkDirAll(name string) error {
	return c.fsys.MkDirAll(name)
}

func (c *Cache) OpenFile(name string) (file.File, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	df, ok := c.files[name]
	if ok {
		return df, nil
	}

	f, err := c.fsys.OpenFile(name)
	if err != nil {
		return nil, err
	}

	if c.files == nil {
		c.files = map[string]*File{}
	}
	df = &File{
		File: f,
		name: name,
		fsys: c}
	c.files[name] = df
	return df, nil
}

func (c *Cache) ReadDir(name string) ([]fs.DirEntry, error) {
	return c.fsys.ReadDir(name)
}

func (c *Cache) Stat(name string) (fs.FileInfo, error) {
	return c.fsys.Stat(name)
}

func (c *Cache) Remove(name string) error {
	if err := c.CloseFile(name); err != nil {
		return err
	}

	return c.fsys.Remove(name)
}

func (c *Cache) Count() int {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return len(c.files)
}

func (c *Cache) Close() error {
	errs := make([]error, 0, len(c.files))
	for n := range c.files {
		errs = append(errs, c.CloseFile(n))
	}
	return errors.Join(errs...)
}

func (c *Cache) CloseFile(name string) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	f, ok := c.files[name]
	if !ok {
		return nil
	}

	if err := f.closeFile(); err != nil {
		return err
	}

	delete(c.files, name)
	return nil
}
