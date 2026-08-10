package fixedv2

import (
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
		fc := New(s.Storage, btree.NewCursor(s.Storage), l)

		fn := link.NewFileName("testfile")
		fnIdx := link.NewFileName("testindex")
		// pid := value.NewPartitionID(0)

		if err := fc.Insert(t.Context(), fn, fnIdx,
			key.NewKey([]byte("abcde")),
			[]byte{1, 2, 3, 4, 5},
		); err != nil {
			t.Fatal(err)
		}
	})
}
