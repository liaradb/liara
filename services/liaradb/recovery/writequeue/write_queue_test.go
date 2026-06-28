package writequeue

import (
	"io"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/record"
)

func TestWriteQueue(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testWriteQueueTestWriteQueue)
}

func testWriteQueueTestWriteQueue(t *testing.T) {
	ps := &testPageStorage{}
	pool := pagepool.New(12, 16, 8)
	wq := New(10, ps, &testFlusher{}, &pool)
	go wq.Run(t.Context())

	wq.Append(t.Context(), record.NewLogSequenceNumber(0), logpage.NewLogPage(128, 16, 8))
	wq.Append(t.Context(), record.NewLogSequenceNumber(1), logpage.NewLogPage(128, 16, 8))
	wq.Replace(t.Context(), logpage.NewLogPage(128, 16, 8))
	if err := wq.ReplaceSync(t.Context(), logpage.NewLogPage(128, 16, 8)); err != nil {
		t.Fatal(err)
	}

	if err := wq.AppendSync(t.Context(), record.NewLogSequenceNumber(2), logpage.NewLogPage(128, 16, 8)); err != nil {
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

func (t *testPageStorage) Append(record.LogSequenceNumber, []byte) error {
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

type testFlusher struct {
}

func (*testFlusher) OnFlush(record.LogSequenceNumber) {}

var _ Flusher = &testFlusher{}
