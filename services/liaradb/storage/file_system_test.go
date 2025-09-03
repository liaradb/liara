package storage

import (
	"path"
	"testing"
)

func TestFileSystem(t *testing.T) {
	p := path.Join(t.TempDir(), "file")
	fs := &FileSystem{}

	if c := fs.Count(); c != 0 {
		t.Errorf("incorrect count, expected: %v, recieved: %v", 0, c)
	}

	f, err := fs.Open(p)
	if err != nil {
		t.Error(err)
	}

	if f == nil {
		t.Error("file should not be nil")
	}

	if c := fs.Count(); c != 1 {
		t.Errorf("incorrect count, expected: %v, recieved: %v", 1, c)
	}

	f, err = fs.Open(p)
	if err != nil {
		t.Error(err)
	}

	if f == nil {
		t.Error("file should not be nil")
	}

	if c := fs.Count(); c != 1 {
		t.Errorf("incorrect count, expected: %v, recieved: %v", 1, c)
	}

	err = fs.Close()
	if err != nil {
		t.Error(err)
	}

	if c := fs.Count(); c != 0 {
		t.Errorf("incorrect count, expected: %v, recieved: %v", 0, c)
	}
}
