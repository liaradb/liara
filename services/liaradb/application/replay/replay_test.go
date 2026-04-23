package replay

import (
	"testing"
	"time"

	"github.com/liaradb/liaradb/collection"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestReplay(t *testing.T) {
	storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
		l := recovery.NewLog(256, 2, s.FSys, "dir")
		r := NewReplay(collection.NewCollections(s.Storage), l)

		if err := l.Open(t.Context()); err != nil {
			t.Fatal(err)
		}

		defer l.Close()

		if err := l.StartWriter(); err != nil {
			t.Fatal(err)
		}

		tid := value.NewTenantID()
		txid := record.NewTransactionID(2)

		if lsn, err := l.Start(t.Context(),
			tid,
			txid,
			time.UnixMicro(1234567890),
		); err != nil {
			t.Error(err)
		} else if lsn != record.NewLogSequenceNumber(1) {
			t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
		}

		if lsn, err := l.Insert(t.Context(),
			tid,
			txid,
			time.UnixMicro(1234567890),
			record.CollectionValue,
			[]byte{1, 2, 3, 4, 5},
		); err != nil {
			t.Error(err)
		} else if lsn != record.NewLogSequenceNumber(2) {
			t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
		}

		if lsn, err := l.Commit(t.Context(),
			tid,
			txid,
			time.UnixMicro(1234567890),
		); err != nil {
			t.Error(err)
		} else if lsn != record.NewLogSequenceNumber(3) {
			t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
		}

		if err := r.Recover(t.Context()); err != nil {
			t.Fatal(err)
		}
	})
}
