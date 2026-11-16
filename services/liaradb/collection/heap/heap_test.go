package heap

import (
	"testing"

	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/logtesting"
	"github.com/liaradb/liaradb/util/testing/should"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestHeap_Get(t *testing.T) {
	storagetesting.SyncTest(t, 100, 1024, func(t *testing.T, s storagetesting.Storage) {
		h := New(s.Storage)

		fn := link.NewFileName("test")
		bid := link.BlockID(0)
		rid := link.NewRecordID(bid, link.SlotID(0))

		should.Panic(t, "should not get value", func() {
			_, _ = h.Get(t.Context(), fn, rid)
		})

		want := NewTestItem(NewTestID())
		rid, err := h.Insert(t.Context(), &logtesting.MockLog{}, fn, want)
		if err != nil {
			t.Fatal(err)
		}

		r, err := h.Get(t.Context(), fn, rid)
		if err != nil {
			t.Fatal(err)
		}

		defer r.Release()

		var item TestItem
		if err := item.Read(r); err != nil {
			t.Fatal(err)
		}

		if item != *want {
			t.Errorf("incorrect result: %v, expected: %v", item, *want)
		}
	})
}

func TestHeap_Delete(t *testing.T) {
	storagetesting.SyncTest(t, 100, 1024, func(t *testing.T, s storagetesting.Storage) {
		h := New(s.Storage)

		fn := link.NewFileName("test")
		bid := link.BlockID(0)
		rid := link.NewRecordID(bid, link.SlotID(0))

		want := NewTestItem(NewTestID())
		rid, err := h.Insert(t.Context(), &logtesting.MockLog{}, fn, want)
		if err != nil {
			t.Fatal(err)
		}

		if err := h.Delete(t.Context(), &logtesting.MockLog{}, fn, rid); err != nil {
			t.Fatal(err)
		}

		should.Panic(t, "should not get value", func() {
			_, _ = h.Get(t.Context(), fn, rid)
		})
	})
}
