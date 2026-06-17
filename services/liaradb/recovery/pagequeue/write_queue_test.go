package pagequeue

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/record"
)

func TestWriteQueue(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testWriteQueueTestWriteQueue)
}

func testWriteQueueTestWriteQueue(t *testing.T) {
	ps := &testPageStorage{}
	wq := newWriteQueue(10, ps)
	wq.Run(t.Context())

	wq.Append(t.Context(), record.NewLogSequenceNumber(0), page.New(128, 16, 8))
	wq.Append(t.Context(), record.NewLogSequenceNumber(1), page.New(128, 16, 8))
	if err := wq.Sync(t.Context(), record.NewLogSequenceNumber(2), page.New(128, 16, 8)); err != nil {
		t.Fatal(err)
	}

	if ps.appendCount != 2 {
		t.Errorf("incorrect append count: %v, expected: %v", ps.appendCount, 2)
	}

	if ps.syncCount != 1 {
		t.Errorf("incorrect sync count: %v, expected: %v", ps.syncCount, 1)
	}
}
