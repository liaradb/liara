package fixedv2

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestFixedCollection(t *testing.T) {
	storagetesting.SyncTest(t, 16, 1024, func(t *testing.T, s storagetesting.Storage) {
		l := log.New(256, 2, 256, 100, s.FSys, "dir")
		if err := l.Run(t.Context()); err != nil {
			t.Fatal(err)
		}

		fc := New(s.Storage, btree.NewCursor(s.Storage), l)

		fn := link.NewFileName("testfile")
		fnIdx := link.NewFileName("testindex")
		// pid := value.NewPartitionID(0)
		k := key.NewKey([]byte("abcde"))
		want := []byte{1, 2, 3, 4, 5}

		if err := fc.Insert(t.Context(), fn, fnIdx, k, want); err != nil {
			t.Fatal(err)
		}

		result, err := fc.Get(t.Context(), fn, fnIdx, k)
		if err != nil {
			t.Fatal(err)
		}

		if !slices.Equal(result, want) {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	})
}
