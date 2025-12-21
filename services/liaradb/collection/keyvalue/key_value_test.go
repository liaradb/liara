package keyvalue

import (
	"testing"

	"github.com/liaradb/liaradb/storage/storagetesting"
)

func TestKeyValue(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	kv := New(s)

	if err := kv.Set(t.Context(), "testfile", "1", []byte("a")); err != nil {
		t.Error(err)
	}

	value, err := kv.Get(t.Context(), "testfile", "1")
	if err != nil {
		t.Error(err)
	}

	want := "a"
	result := string(value)
	if result != want {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}
