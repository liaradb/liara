package disk

import (
	"path"
	"testing"
)

func TestFileSystem(t *testing.T) {
	t.Parallel()

	p := path.Join(t.TempDir(), "file")
	fs := NewFileSystem(&DiskFileSystem{})

	if c := fs.Count(); c != 0 {
		t.Errorf("incorrect count, expected: %v, recieved: %v", 0, c)
	}

	if f, err := fs.OpenFile(p); err != nil {
		t.Error(err)
	} else if f == nil {
		t.Error("file should not be nil")
	}

	if c := fs.Count(); c != 1 {
		t.Errorf("incorrect count, expected: %v, recieved: %v", 1, c)
	}

	if f, err := fs.OpenFile(p); err != nil {
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

func TestFileSystem_CloseFile(t *testing.T) {
	t.Parallel()

	t.Run("should close", func(t *testing.T) {
		t.Parallel()

		p := path.Join(t.TempDir(), "file")
		fs := NewFileSystem(&DiskFileSystem{})
		if f, err := fs.OpenFile(p); err != nil {
			t.Error(err)
		} else if f == nil {
			t.Error("file should not be nil")
		}

		if err := fs.CloseFile(p); err != nil {
			t.Error(err)
		}
	})

	t.Run("should close individual files", func(t *testing.T) {
		t.Parallel()

		p := path.Join(t.TempDir(), "file")
		fs := NewFileSystem(&DiskFileSystem{})

		f, err := fs.OpenFile(p)
		if err != nil {
			t.Fatal(err)
		} else if f == nil {
			t.Fatal("file should not be nil")
		}

		if err := f.Close(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Should noop if file not opened", func(t *testing.T) {
		t.Parallel()

		fs := NewFileSystem(&DiskFileSystem{})

		if err := fs.CloseFile("other"); err != nil {
			t.Error(err)
		}
	})

	t.Run("should noop if file already closed", func(t *testing.T) {
		t.Parallel()

		p := path.Join(t.TempDir(), "file")
		fs := NewFileSystem(&DiskFileSystem{})

		f, err := fs.OpenFile(p)
		if err != nil {
			t.Error(err)
		} else if f == nil {
			t.Error("file should not be nil")
		}

		if err := f.Close(); err != nil {
			t.Fatal(err)
		}

		if err := fs.CloseFile(p); err != nil {
			t.Error(err)
		}
	})
}
