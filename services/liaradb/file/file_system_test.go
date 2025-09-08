package file

import (
	"path"
	"testing"
)

func TestFileSystem(t *testing.T) {
	t.Parallel()

	p := path.Join(t.TempDir(), "file")
	fs := &FileSystem{}

	if c := fs.Count(); c != 0 {
		t.Errorf("incorrect count, expected: %v, recieved: %v", 0, c)
	}

	if f, err := fs.Open(p); err != nil {
		t.Error(err)
	} else if f == nil {
		t.Error("file should not be nil")
	}

	if c := fs.Count(); c != 1 {
		t.Errorf("incorrect count, expected: %v, recieved: %v", 1, c)
	}

	if f, err := fs.Open(p); err != nil {
		t.Error(err)
	} else if f == nil {
		t.Error("file should not be nil")
	}

	if c := fs.Count(); c != 1 {
		t.Errorf("incorrect count, expected: %v, recieved: %v", 1, c)
	}

	if err := fs.Close(); err != nil {
		t.Error(err)
	}

	if c := fs.Count(); c != 0 {
		t.Errorf("incorrect count, expected: %v, recieved: %v", 0, c)
	}
}
