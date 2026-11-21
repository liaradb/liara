package tip

import (
	"context"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
	"github.com/liaradb/liaradb/util/testing/txtesting"
)

const (
	pageSize       = 78
	largePageSize  = 256
	writeQueueSize = 100
)

func TestTip(t *testing.T) {
	storagetesting.SyncTest(t, 16, pageSize, func(t *testing.T, st storagetesting.Storage) {
		tip := NewTip(st.Storage, link.NewFileName("fn"))
		want := 128
		s, err := tip.Record(t.Context(), want)
		if err != nil {
			t.Fatal(err)
		}

		if l := s.Length(); l != want {
			t.Errorf("incorrect length: %v, expected: %v", l, want)
		}

		tx := txtesting.Transaction{}
		if err := s.CommitAppend(t.Context(), &tx); err != nil {
			t.Fatal(err)
		}

		s.Release()

		size := 0
		for _, r := range tx.Records() {
			size += len(r.Data())
		}
		if size != want {
			t.Errorf("incorrect length: %v, expected: %v", tx, want)
		}
		if c := len(tx.Records()); c != 3 {
			t.Errorf("incorrect complete count: %v, expected: %v", c, 3)
		}

		synctest.Wait()
	})
}

func TestTip__Error(t *testing.T) {
	storagetesting.SyncTest(t, 16, pageSize, func(t *testing.T, st storagetesting.Storage) {
		tip := NewTip(st.Storage, link.NewFileName("fn"))

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		_, err := tip.Record(ctx, 128)
		if err == nil {
			t.Error("should return error")
		}
	})
}
