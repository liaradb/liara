package pagequeue

import (
	"sync"
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
	wg := sync.WaitGroup{}

	ps := &testPageStorage{}
	wq := newWriteQueue(10, ps, &testFlusher{
		onFlush: func(lsn record.LogSequenceNumber) {
			if v := lsn.Value(); v != 1 {
				t.Errorf("incorrect lsn: %v, expected: %v", v, 1)
			}
			wg.Done()
		},
	})
	wq.Run(t.Context())

	wg.Add(1)
	wq.Append(t.Context(), record.NewLogSequenceNumber(0), page.New(128, 16, 8))
	wq.Append(t.Context(), record.NewLogSequenceNumber(1), page.New(128, 16, 8))
	wg.Wait()

	if ps.appendCount != 2 {
		t.Errorf("incorrect append count: %v, expected: %v", ps.appendCount, 2)
	}

	if ps.syncCount != 0 {
		t.Errorf("incorrect sync count: %v, expected: %v", ps.syncCount, 0)
	}
}

type testFlusher struct {
	onFlush func(record.LogSequenceNumber)
}

func (t *testFlusher) OnError(error) bool {
	return true
}

func (t *testFlusher) OnFlush(lsn record.LogSequenceNumber) {
	t.onFlush(lsn)
}

var _ Flusher = &testFlusher{}
