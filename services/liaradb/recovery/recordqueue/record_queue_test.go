package recordqueue

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/pagepool"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/recovery/span"
	"github.com/liaradb/liaradb/transaction/record"
	"github.com/liaradb/liaradb/util/testing/filetesting"
)

func TestRecordQueue_Default(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testRecordQueue_Default)
}

func testRecordQueue_Default(t *testing.T) {
	rq := createRecordQueueStart(t, 320, 3, 320)

	testPosition(t, rq, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(0))
}

func TestRecordQueue_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testRecordQueue_Append)
}

func testRecordQueue_Append(t *testing.T) {
	ctx := t.Context()

	rq := createRecordQueueStart(t, 320, 3, 320)
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	if lsn, err := rq.Append(ctx, record.New(
		value.NewTenantID(),
		record.NewTransactionID(2),
		record.NewTime(time.UnixMicro(1234567890)),
		record.ActionUpdate,
		record.CollectionValue,
		data,
		reverse,
	)); err != nil {
		t.Error(err)
	} else if lsn != logpage.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, rq, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(1))
}

func TestRecordQueue_Append__Large(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testRecordQueue_Append__Large)
}

func testRecordQueue_Append__Large(t *testing.T) {
	ctx := t.Context()

	rq := createRecordQueueStart(t, 320, 3, 320)
	var data = make([]byte, 0, 1024)
	for i := range 1024 {
		data = append(data, byte(i%255))
	}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	if _, err := rq.Append(ctx, record.New(
		value.NewTenantID(),
		record.NewTransactionID(2),
		record.NewTime(time.UnixMicro(1234567890)),
		record.ActionUpdate,
		record.CollectionValue,
		data,
		reverse,
	)); err != raw.ErrInsufficientSpace {
		t.Errorf("should return %v", raw.ErrInsufficientSpace)
	}

	testPosition(t, rq, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(0))
}

func TestRecordQueue_Flush(t *testing.T) {
	t.Parallel()

	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	t.Run("should flush", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			rq := createRecordQueueStart(t, 320, 3, 320)
			tid := value.NewTenantID()

			if _, err := rq.Append(ctx, record.New(
				tid,
				record.NewTransactionID(2),
				record.NewTime(time.UnixMicro(1234567890)),
				record.ActionUpdate,
				record.CollectionValue,
				data,
				reverse,
			)); err != nil {
				t.Error(err)
			}

			testPosition(t, rq, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(1))

			if _, err := rq.Append(ctx, record.New(
				tid,
				record.NewTransactionID(2),
				record.NewTime(time.UnixMicro(1234567890)),
				record.ActionUpdate,
				record.CollectionValue,
				data,
				reverse,
			)); err != nil {
				t.Error(err)
			}

			testPosition(t, rq, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(2))

			if err := rq.Flush(t.Context()); err != nil {
				t.Error(err)
			}

			testPosition(t, rq, logpage.NewLogSequenceNumber(2), logpage.NewLogSequenceNumber(2))
		})
	})

	t.Run("should not flush beyond HighWater", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			rq := createRecordQueueStart(t, 320, 3, 320)
			tid := value.NewTenantID()

			if _, err := rq.Append(ctx, record.New(
				tid,
				record.NewTransactionID(2),
				record.NewTime(time.UnixMicro(1234567890)),
				record.ActionUpdate,
				record.CollectionValue,
				data,
				reverse,
			)); err != nil {
				t.Error(err)
			}

			if _, err := rq.Append(ctx, record.New(
				tid,
				record.NewTransactionID(2),
				record.NewTime(time.UnixMicro(1234567890)),
				record.ActionUpdate,
				record.CollectionValue,
				data,
				reverse,
			)); err != nil {
				t.Error(err)
			}

			if err := rq.Flush(t.Context()); err != nil {
				t.Error(err)
			}

			testPosition(t, rq, logpage.NewLogSequenceNumber(2), logpage.NewLogSequenceNumber(2))
		})
	})

	t.Run("should write to multiple pages", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			rq := createRecordQueueStart(t, 352, 4, 352)
			tid := value.NewTenantID()
			count := 14
			for range count {
				if _, err := rq.Append(ctx, record.New(
					tid,
					record.NewTransactionID(2),
					record.NewTime(time.UnixMicro(1234567890)),
					record.ActionUpdate,
					record.CollectionValue,
					data,
					reverse,
				)); err != nil {
					t.Fatal(err)
				}
			}

			if err := rq.Flush(t.Context()); err != nil {
				t.Error(err)
			}

			// if p := l.PageID(); p != 3 {
			// 	t.Errorf("incorrect value: %v, expected: %v", p, 3)
			// }
		})
	})

	t.Run("should return error if appending beyond maximum", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			rq := createRecordQueueStart(t, 48, 1, 2)

			if _, err := rq.Append(ctx, record.New(
				value.NewTenantID(),
				record.NewTransactionID(2),
				record.NewTime(time.UnixMicro(1234567890)),
				record.ActionUpdate,
				record.CollectionValue,
				data,
				reverse,
			)); err == nil {
				t.Fatal("should return error")
			}
		})
	})

	t.Run("should write after flushing", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			rq := createRecordQueueStart(t, 320, 3, 320)
			tid := value.NewTenantID()

			if _, err := rq.Append(ctx, record.New(
				tid,
				record.NewTransactionID(2),
				record.NewTime(time.UnixMicro(1234567890)),
				record.ActionUpdate,
				record.CollectionValue,
				data,
				reverse,
			)); err != nil {
				t.Error(err)
			}

			if err := rq.Flush(t.Context()); err != nil {
				t.Error(err)
			}

			if _, err := rq.Append(ctx, record.New(
				tid,
				record.NewTransactionID(2),
				record.NewTime(time.UnixMicro(1234567890)),
				record.ActionUpdate,
				record.CollectionValue,
				data,
				reverse,
			)); err != nil {
				t.Error(err)
			}

			if err := rq.Flush(t.Context()); err != nil {
				t.Error(err)
			}

			testPosition(t, rq, logpage.NewLogSequenceNumber(2), logpage.NewLogSequenceNumber(2))
		})
	})
}

func createRecordQueueStart(t *testing.T,
	pageSize int64,
	segmentSize segment.PageID,
	recordSize int64,
) *RecordQueue[*record.Record] {
	t.Helper()

	fsys, dir := createFiles()
	rq := createRecordQueue(t, pageSize, segmentSize, recordSize, fsys, dir)
	if err := rq.Init(logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(0)); err != nil {
		t.Fatal(err)
	}

	return rq
}

func createRecordQueue(t *testing.T,
	pageSize int64,
	segmentSize segment.PageID,
	recordSize int64,
	fsys filecache.FileSystem,
	dir string,
) *RecordQueue[*record.Record] {
	t.Helper()

	sl := segment.NewList(fsys, dir, pageSize, segmentSize)
	// TODO: Fix this cast
	pl := pagepool.New(int(pageSize), span.FragmentHeaderSize)

	rq := New[*record.Record](pageSize, recordSize, 100, sl, pl)
	if err := rq.Run(t.Context()); err != nil {
		t.Fatal(err)
	}

	cleanupRecordQueue(t, rq)

	return rq
}

func cleanupRecordQueue(t *testing.T, l *RecordQueue[*record.Record]) {
	t.Helper()

	t.Cleanup(func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	})
}

func createFiles() (filecache.FileSystem, string) {
	// return &disk.FileSystem{}, t.TempDir()
	return filetesting.New(nil), "."
}

func testPosition(t *testing.T, l *RecordQueue[*record.Record], lw, hw logpage.LogSequenceNumber) {
	t.Helper()

	if h := l.HighWater(); h != hw {
		t.Errorf("incorrect high water: %v, expected: %v", h, hw)
	}

	if l := l.LowWater(); l != lw {
		t.Errorf("incorrect low water: %v, expected: %v", l, lw)
	}
}
