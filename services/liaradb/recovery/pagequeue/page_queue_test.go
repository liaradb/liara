package pagequeue

import (
	"io"
	"reflect"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/pagequeue/writequeue"
	"github.com/liaradb/liaradb/recovery/span"
)

func TestPageQueue(t *testing.T) {
	t.Parallel()

	lsn := logpage.NewLogSequenceNumber(1)

	t.Run("should run", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			rc := newTestRecord(100)
			rc.SetLogSequenceNumber(lsn)

			ps := &testPageStorage{}
			pq := New[*testRecord](ps, pagepool.New(largePageSize, span.FragmentHeaderSize), writeQueueSize)

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

			rc2 := newTestRecord(100)
			if err := rc2.Read(s); err != nil {
				t.Fatal(err)
			}

			if i := rc2.LogSequenceNumber(); i != lsn {
				t.Errorf("incorrect log sequence number: %v, expected: %v", i, lsn)
			}

			if !reflect.DeepEqual(rc2, rc) {
				t.Errorf("incorrect record: %v, expected: %v", rc2, rc)
			}
		})
	})

	t.Run("should run next", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			rc := newTestRecord(100)
			rc.SetLogSequenceNumber(lsn)

			ps := &testPageStorage{}
			pq := New[*testRecord](ps, pagepool.New(largePageSize, span.FragmentHeaderSize), writeQueueSize)

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
			rc := newTestRecord(100)
			rc.SetLogSequenceNumber(lsn)

			ps := &testPageStorage{}
			pq := New[*testRecord](ps, pagepool.New(largePageSize, span.FragmentHeaderSize), writeQueueSize)
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
			rc := newTestRecord(100)
			rc.SetLogSequenceNumber(lsn)

			ps := &testPageStorage{}
			pq := New[*testRecord](ps, pagepool.New(largePageSize, span.FragmentHeaderSize), writeQueueSize)
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
			rc := newTestRecord(80)
			rc.SetLogSequenceNumber(lsn)

			ps := &testPageStorage{}
			pq := New[*testRecord](ps, pagepool.New(largePageSize, span.FragmentHeaderSize), writeQueueSize)
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
			rc := newTestRecord(100)
			rc.SetLogSequenceNumber(lsn)

			ps := &testPageStorage{}
			pq := New[*testRecord](ps, pagepool.New(largePageSize, span.FragmentHeaderSize), writeQueueSize)
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
			rc := newTestRecord(100)
			rc.SetLogSequenceNumber(lsn)

			ps := &testPageStorage{}
			pq := New[*testRecord](ps, pagepool.New(largePageSize, span.FragmentHeaderSize), writeQueueSize)
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

func (t *testPageStorage) Append(logpage.LogSequenceNumber, []byte) error {
	if t.errorOnAppend {
		return writequeue.ErrUnableToAppend
	}

	t.appendCount++
	return nil
}

func (t *testPageStorage) Init([]byte) error {
	return nil
}

func (t *testPageStorage) Replace([]byte) error {
	if t.errorOnSync {
		return io.ErrShortWrite
	}

	t.syncCount++
	return nil
}

var _ writequeue.PageStorage = &testPageStorage{}
