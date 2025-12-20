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
}
