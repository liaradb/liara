package writequeue

import (
	"io"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
)

func TestWriteQueue(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testWriteQueueTestWriteQueue)
}

func testWriteQueueTestWriteQueue(t *testing.T) {
	ps := &testPageStorage{}
	pool := pagepool.New(12, 8)
	wq := New(10, ps, pool)
	go wq.Run(t.Context())

	wq.Append(t.Context(), logpage.NewLogSequenceNumber(0), logpage.New(128, 8))
	wq.Append(t.Context(), logpage.NewLogSequenceNumber(1), logpage.New(128, 8))
	wq.Replace(t.Context(), logpage.New(128, 8))
	if err := wq.ReplaceSync(t.Context(), logpage.New(128, 8)); err != nil {
		t.Fatal(err)
	}

	if err := wq.AppendSync(t.Context(), logpage.NewLogSequenceNumber(2), logpage.New(128, 8)); err != nil {
		t.Fatal(err)
	}

	if ps.appendCount != 3 {
		t.Errorf("incorrect append count: %v, expected: %v", ps.appendCount, 2)
	}

	if ps.syncCount != 2 {
		t.Errorf("incorrect sync count: %v, expected: %v", ps.syncCount, 2)
	}
}

type testPageStorage struct {
	errorOnAppend bool
	errorOnSync   bool
	syncCount     int
	appendCount   int
}

func (t *testPageStorage) Append(logpage.LogSequenceNumber, []byte) error {
	if t.errorOnAppend {
		return ErrUnableToAppend
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

var _ PageStorage = &testPageStorage{}
