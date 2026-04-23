package filetesting

import (
	"io"
	"slices"
	"testing"
)

func TestFileSystem(t *testing.T) {
	t.Parallel()

	fc := New(nil)
	if fc == nil {
		t.Fatal("should return value")
	}

	if fsys, ok := fc.FSYS().(*FileSystem); !ok {
		t.Fatal("incorrect type")
	} else if fsys == nil {
		t.Fatal("should return value")
	}

	f, err := fc.OpenFile("file")
	if err != nil {
		t.Fatal(err)
	}

	want := []byte("abcde")
	if n, err := f.Write(want); err != nil {
		t.Fatal(err)
	} else if n != 5 {
		t.Errorf("incorrect length: %v, expected: %v", n, 5)
	}

	if n, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	} else if n != 0 {
		t.Errorf("incorrect length: %v, expected: %v", n, 0)
	}

	result := make([]byte, 5)
	if n, err := f.Read(result); err != nil {
		t.Fatal(err)
	} else if n != 5 {
		t.Errorf("incorrect length: %v, expected: %v", n, 5)
	}

	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := f.Seek(10, io.SeekStart); err == nil {
		t.Error("should return error")
	}

	f, err = fc.OpenFile("file")
	if err != nil {
		t.Fatal(err)
	}

	result = make([]byte, 5)
	if n, err := f.Read(result); err != nil {
		t.Fatal(err)
	} else if n != 5 {
		t.Errorf("incorrect length: %v, expected: %v", n, 5)
	}

	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}
