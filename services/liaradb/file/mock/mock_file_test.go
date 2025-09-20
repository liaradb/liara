package mock

import (
	"encoding/binary"
	"testing"
)

func TestFile(t *testing.T) {
	t.Parallel()

	f := NewMockFile("file")
	f.Open()

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

	s, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if s.Name() != "file" {
		t.Error("wrong name")
	}
	if s.Size() != 108 {
		t.Error("wrong size")
	}
}
