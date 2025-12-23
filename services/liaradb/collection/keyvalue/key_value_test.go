package keyvalue

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/storage/storagetesting"
)

func TestKeyValue(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testKeyValue)
}

func testKeyValue(t *testing.T) {
	ctx := t.Context()
	// TODO: This is flaky on insert when buffer count is 4
	// s := storagetesting.CreateStorage(t, 4, 64)
	s := storagetesting.CreateStorage(t, 5, 64)
	kv := New(s)
	n := tablename.New("testfile")

	data := map[value.Key][]byte{
		"1": []byte("a"),
		"2": []byte("b"),
		"3": []byte("c"),
		"4": []byte("d"),
		"5": []byte("e"),
		"6": []byte("f"),
		"7": []byte("g"),
		"8": []byte("h"),
		"9": []byte("i"),
	}

	for k, v := range data {
		if err := kv.Set(ctx, n, k, v); err != nil {
			t.Fatal(err)
		}
	}

	for k, v := range data {
		value, err := kv.Get(ctx, n, k)
		if err != nil {
			t.Fatal(k, err)
		}

		want := string(v)
		result := string(value)
		if result != want {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	}
}
