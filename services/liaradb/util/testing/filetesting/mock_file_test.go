package filetesting

import (
	"encoding/binary"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"slices"
	"testing"
	"time"
)

func TestFile(t *testing.T) {
	t.Parallel()

	now := time.Now()
	f := NewMockFile("file", 0, now)
	f.Open()

	s, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if mod := s.ModTime(); !mod.Equal(now) {
		t.Errorf("incorrect mod time: %v, expected: %v", mod, now)
	}

	wb := make([]byte, 8)
	binary.LittleEndian.PutUint64(wb, 12345)
	count, err := f.WriteAt(wb, 100)
	if count != 8 {
		t.Fatal("wrong count")
	}
	if err != nil {
		t.Fatal(err)
	}

	rb := make([]byte, 8)
	count, err = f.ReadAt(rb, 100)
	if count != 8 {
		t.Fatal("wrong count")
	}
	if err != nil {
		t.Fatal(err)
	}

	value := binary.LittleEndian.Uint64(rb)
	if value != 12345 {
		t.Fatal("wrong value")
	}

	s, err = f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if n := s.Name(); n != "file" {
		t.Errorf("incorrect name: %v, expected: %v", n, "file")
	}
	if size := s.Size(); size != 108 {
		t.Errorf("incorrect size: %v, expected: %v", size, 108)
	}
	if m := s.Mode(); m != fs.ModeAppend {
		t.Errorf("incorrect mode: %v, expected: %v", m, fs.ModeAppend)
	}
}

func TestFile_Write(t *testing.T) {
	t.Parallel()

	f := NewMockFile("file", 0, time.Time{})
	f.Open()

	data0 := []byte{1, 2}
	data1 := []byte{3, 4, 5}

	if p, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	} else if p != 0 {
		t.Fatalf("incorrect position: %v, expected: %v", p, 0)
	}

	if p, err := f.Write(data0); err != nil {
		t.Fatal(err)
	} else if p != 2 {
		t.Fatalf("incorrect position: %v, expected: %v", p, 2)
	}

	if p, err := f.Write(data1); err != nil {
		t.Fatal(err)
	} else if p != 3 {
		t.Fatalf("incorrect position: %v, expected: %v", p, 3)
	}

	if p, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	} else if p != 0 {
		t.Fatalf("incorrect position: %v, expected: %v", p, 0)
	}

	result0 := make([]byte, 2)
	if p, err := f.Read(result0); err != nil {
		t.Fatal(err)
	} else if p != 2 {
		t.Fatalf("incorrect position: %v, expected: %v", p, 2)
	} else if !slices.Equal(result0, data0) {
		t.Errorf("incorrect result: %v, expected: %v", result0, data0)
	}

	result1 := make([]byte, 3)
	if p, err := f.Read(result1); err != nil {
		t.Fatal(err)
	} else if p != 3 {
		t.Fatalf("incorrect position: %v, expected: %v", p, 3)
	} else if !slices.Equal(result1, data1) {
		t.Errorf("incorrect result: %v, expected: %v", result1, data1)
	}

	if p, err := f.Seek(2, io.SeekStart); err != nil {
		t.Fatal(err)
	} else if p != 2 {
		t.Fatalf("incorrect position: %v, expected: %v", p, 2)
	}

	result2 := make([]byte, 3)
	if p, err := f.Read(result2); err != nil {
		t.Fatal(err)
	} else if p != 3 {
		t.Fatalf("incorrect position: %v, expected: %v", p, 3)
	} else if !slices.Equal(result2, data1) {
		t.Errorf("incorrect result: %v, expected: %v", result2, data1)
	}
}

func TestFile_Stat__Closed(t *testing.T) {
	f := NewMockFile("file", 0, time.Time{})
	f.Open()

	_, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if err = f.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = f.Stat()
	if !errors.Is(err, fs.ErrClosed) {
		t.Fatal(err)
	}
}

// This test is to verify matched behavior for mock file
func TestFile_Stat__Closed_Verification(t *testing.T) {
	dir := t.TempDir()
	f, err := os.OpenFile(path.Join(dir, "file"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if err = f.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = f.Stat()
	if !errors.Is(err, fs.ErrClosed) {
		t.Fatal(err)
	}
}
