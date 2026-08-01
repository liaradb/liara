package collection

import (
	"testing"

	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestCollections(t *testing.T) {
	storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
		l := log.New(256, 2, 256, 100, s.FSys, "dir")
		c := NewCollections(s.Storage, l)
		if c == nil {
			t.Error("should have value")
		}
		if c.EventLog == nil {
			t.Error("should have value")
		}
		if c.Idempotency == nil {
			t.Error("should have value")
		}
		if c.Outbox == nil {
			t.Error("should have value")
		}
	})
}
