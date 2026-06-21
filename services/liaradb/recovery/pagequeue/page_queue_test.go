package pagequeue

import (
	"io"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/span"
	"github.com/liaradb/liaradb/recovery/writequeue"
)

func TestPageQueue(t *testing.T) {
	t.Parallel()

	lsn := record.NewLogSequenceNumber(1)
	tid := value.NewTenantID()
	txid := record.NewTransactionID(2)
	now := record.NewTime(time.UnixMicro(1234567890))
	action := record.ActionInsert
	collection := record.CollectionEvent
	data := []byte("abcdef")
	reverse := []byte("fghij")

	t.Run("should run", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			rc := record.New(lsn, tid, txid, now, action, collection, data, reverse)

			ps := &testPageStorage{}
			pq := New(ps, largePageSize, headerSize)

			if err := pq.Append(t.Context(), rc); err != nil {
				t.Fatal(err)
			}

			// if c := pq.Count(); c != 1 {
			// 	t.Fatalf("incorrect count: %v, expected: %v", c, 1)
			// }

			s := span.Span{}
			for h, d := range pq.current.Slots() {
				s.Append(h, d)
				break
			}
			s.InitIndexes()

			rc2 := &record.Record{}
			if err := rc2.Read(s); err != nil {
				t.Fatal(err)
			}

			if i := rc2.LogSequenceNumber(); i != lsn {
				t.Errorf("incorrect log sequence number: %v, expected: %v", i, lsn)
			}

			if i := rc2.TenantID(); i != tid {
				t.Errorf("incorrect tenant id: %v, expected: %v", i, tid)
			}

			if i := rc2.TransactionID(); i != txid {
				t.Errorf("incorrect transaction id: %v, expected: %v", i, txid)
			}

			if i := rc2.Time(); i != now {
				t.Errorf("incorrect time: %v, expected: %v", i, now)
			}

			if i := rc2.Action(); i != action {
				t.Errorf("incorrect action: %v, expected: %v", i, action)
			}

			if i := rc2.Collection(); i != collection {
				t.Errorf("incorrect collection: %v, expected: %v", i, collection)
			}

			if i := rc2.Data(); !slices.Equal(i, data) {
				t.Errorf("incorrect data: %v, expected: %v", i, data)
			}

			if i := rc2.Reverse(); !slices.Equal(i, reverse) {
				t.Errorf("incorrect reverse: %v, expected: %v", i, reverse)
			}

			if i := rc2.IsCheckpoint(); i != (action == record.ActionCheckpoint) {
				t.Errorf("incorrect is checkpoint: %v, expected: %v", i, action == record.ActionCheckpoint)
			}
		})
	})

	t.Run("should run next", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			rc := record.New(lsn, tid, txid, now, action, collection, data, reverse)

			ps := &testPageStorage{}
			pq := New(ps, largePageSize, headerSize)

			if err := pq.Append(t.Context(), rc); err != nil {
				t.Fatal(err)
			}

			if err := pq.Append(t.Context(), rc); err != nil {
				t.Fatal(err)
			}

			if err := pq.Append(t.Context(), rc); err != nil {
				t.Fatal(err)
			}

			// if c := pq.Count(); c != 2 {
			// 	t.Fatalf("incorrect count: %v, expected: %v", c, 2)
			// }
		})
	})

	t.Run("should flush one", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			rc := record.New(lsn, tid, txid, now, action, collection, data, reverse)

			ps := &testPageStorage{}
			pq := New(ps, largePageSize, headerSize)
			go pq.Run(t.Context())

			if err := pq.Append(t.Context(), rc); err != nil {
				t.Fatal(err)
			}

			if err := pq.Flush(t.Context()); err != nil {
				t.Error(err)
			}

			synctest.Wait()

			if ps.syncCount != 0 {
				t.Errorf("incorrect sync count: %v, expected: %v", ps.syncCount, 0)
			}

			if ps.appendCount != 1 {
				t.Errorf("incorrect append count: %v, expected: %v", ps.appendCount, 1)
			}
		})
	})

	t.Run("should handle error on flush one", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			rc := record.New(lsn, tid, txid, now, action, collection, data, reverse)

			ps := &testPageStorage{}
			pq := New(ps, largePageSize, headerSize)
			go pq.Run(t.Context())

			if err := pq.Append(t.Context(), rc); err != nil {
				t.Fatal(err)
			}

			ps.errorOnAppend = true

			if err := pq.Flush(t.Context()); err == nil {
				t.Error("should return error")
			}
		})
	})

	t.Run("should flush many", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			rc := record.New(lsn, tid, txid, now, action, collection, data, reverse)

			ps := &testPageStorage{}
			pq := New(ps, largePageSize, headerSize)
			go pq.Run(t.Context())

			for range 6 {
				if err := pq.Append(t.Context(), rc); err != nil {
					t.Fatal(err)
				}
			}

			if err := pq.Flush(t.Context()); err != nil {
				t.Error(err)
			}

			synctest.Wait()

			if ps.syncCount != 0 {
				t.Errorf("incorrect sync count: %v, expected: %v", ps.syncCount, 0)
			}

			if ps.appendCount != 3 {
				t.Errorf("incorrect append count: %v, expected: %v", ps.appendCount, 3)
			}
		})
	})

	t.Run("should handle error on flush many", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			rc := record.New(lsn, tid, txid, now, action, collection, data, reverse)

			ps := &testPageStorage{}
			pq := New(ps, largePageSize, headerSize)
			go pq.Run(t.Context())

			for range 6 {
				if err := pq.Append(t.Context(), rc); err != nil {
					t.Fatal(err)
				}
			}

			if err := pq.Flush(t.Context()); err != nil {
				t.Error(err)
			}

			ps.errorOnSync = true
			if err := pq.Flush(t.Context()); err == nil {
				t.Error("should return error")
			}

			ps.errorOnSync = false
			ps.errorOnAppend = true
			if err := pq.Flush(t.Context()); err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("should clear", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			rc := record.New(lsn, tid, txid, now, action, collection, data, reverse)

			ps := &testPageStorage{}
			pq := New(ps, largePageSize, headerSize)
			go pq.Run(t.Context())

			for range 6 {
				if err := pq.Append(t.Context(), rc); err != nil {
					t.Fatal(err)
				}
			}

			// pq.Clear()
			// if c := pq.Count(); c != 1 {
			// 	t.Errorf("incorrect count: %v, expected: %v", c, 1)
			// }
		})
	})
}

type testPageStorage struct {
	errorOnAppend bool
	errorOnSync   bool
	syncCount     int
	appendCount   int
}

func (t *testPageStorage) Append(record.LogSequenceNumber, []byte) error {
	if t.errorOnAppend {
		return writequeue.ErrUnableToAppend
	}

	t.appendCount++
	return nil
}

func (t *testPageStorage) Init([]byte) error {
	return nil
}

func (t *testPageStorage) Sync([]byte) error {
	if t.errorOnSync {
		return io.ErrShortWrite
	}

	t.syncCount++
	return nil
}

var _ PageStorage = &testPageStorage{}
