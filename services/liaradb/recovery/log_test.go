package recovery

import (
	"reflect"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/util/testing/filetesting"
)

func TestLog_Default(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Default)
}

func testLog_Default(t *testing.T) {
	l := createLogStart(t, 320, 3)

	testPosition(t, l, record.NewLogSequenceNumber(0), record.NewLogSequenceNumber(0))
}

func TestLog_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Append)
}

func testLog_Append(t *testing.T) {
	ctx := t.Context()

	l := createLogStart(t, 320, 3)
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	if lsn, err := l.Update(ctx,
		value.NewTenantID(),
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
		data,
		reverse,
	); err != nil {
		t.Error(err)
	} else if lsn != record.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, record.NewLogSequenceNumber(0), record.NewLogSequenceNumber(1))
}

// TODO: Should not create next Segment if cannot fit
func TestLog_Append__Large(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Append__Large)
}

func testLog_Append__Large(t *testing.T) {
	ctx := t.Context()

	l := createLogStart(t, 320, 3)
	var data = make([]byte, 0, 1024)
	for i := range 1024 {
		data = append(data, byte(i%255))
	}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	if _, err := l.Update(ctx,
		value.NewTenantID(),
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
		data,
		reverse,
	); err != raw.ErrInsufficientSpace {
		t.Errorf("should return %v", raw.ErrInsufficientSpace)
	}

	testPosition(t, l, record.NewLogSequenceNumber(0), record.NewLogSequenceNumber(0))
}

func TestLog_Flush(t *testing.T) {
	t.Parallel()

	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	runTest(t, "should flush", func(t *testing.T) {
		ctx := t.Context()

		l := createLogStart(t, 320, 3)

		tid := value.NewTenantID()

		_, err := l.Update(ctx,
			tid,
			record.NewTransactionID(2),
			time.UnixMicro(1234567890),
			data,
			reverse)
		if err != nil {
			t.Error(err)
		}

		testPosition(t, l, record.NewLogSequenceNumber(0), record.NewLogSequenceNumber(1))

		if _, err = l.Update(ctx,
			tid,
			record.NewTransactionID(2),
			time.UnixMicro(1234567890),
			data,
			reverse,
		); err != nil {
			t.Error(err)
		}

		testPosition(t, l, record.NewLogSequenceNumber(0), record.NewLogSequenceNumber(2))

		if err := l.Flush(ctx); err != nil {
			t.Error(err)
		}

		testPosition(t, l, record.NewLogSequenceNumber(2), record.NewLogSequenceNumber(2))
	})

	runTest(t, "should not flush beyond HighWater", func(t *testing.T) {
		ctx := t.Context()

		l := createLogStart(t, 320, 3)
		tid := value.NewTenantID()

		if _, err := l.Update(ctx,
			tid,
			record.NewTransactionID(2),
			time.UnixMicro(1234567890),
			data,
			reverse,
		); err != nil {
			t.Error(err)
		}

		if _, err := l.Update(ctx,
			tid,
			record.NewTransactionID(2),
			time.UnixMicro(1234567890),
			data,
			reverse,
		); err != nil {
			t.Error(err)
		}

		if err := l.Flush(ctx); err != nil {
			t.Error(err)
		}

		testPosition(t, l, record.NewLogSequenceNumber(2), record.NewLogSequenceNumber(2))
	})

	runTest(t, "should write to multiple pages", func(t *testing.T) {
		ctx := t.Context()

		l := createLogStart(t, 344, 4)
		tid := value.NewTenantID()

		count := 14

		for range count {
			if _, err := l.Update(ctx,
				tid,
				record.NewTransactionID(2),
				time.UnixMicro(1234567890),
				data,
				reverse,
			); err != nil {
				t.Fatal(err)
			}
		}

		if err := l.Flush(ctx); err != nil {
			t.Fatal(err)
		}

		if p := l.PageID(); p != 3 {
			t.Errorf("incorrect value: %v, expected: %v", p, 3)
		}
	})

	runTest(t, "should return error if appending beyond maximum", func(t *testing.T) {
		ctx := t.Context()

		l := createLogStart(t, 32, 1)

		if _, err := l.Update(ctx,
			value.NewTenantID(),
			record.NewTransactionID(2),
			time.UnixMicro(1234567890),
			data,
			reverse,
		); err != raw.ErrInsufficientSpace {
			t.Fatal("should return error")
		}
	})

	runTest(t, "should write after flushing", func(t *testing.T) {
		ctx := t.Context()

		l := createLogStart(t, 320, 3)
		tid := value.NewTenantID()

		if _, err := l.Update(ctx,
			tid,
			record.NewTransactionID(2),
			time.UnixMicro(1234567890),
			data,
			reverse); err != nil {
			t.Error(err)
		}

		if err := l.Flush(ctx); err != nil {
			t.Error(err)
		}

		if _, err := l.Update(ctx,
			tid,
			record.NewTransactionID(2),
			time.UnixMicro(1234567890),
			data,
			reverse); err != nil {
			t.Error(err)
		}

		if err := l.Flush(ctx); err != nil {
			t.Error(err)
		}

		testPosition(t, l, record.NewLogSequenceNumber(2), record.NewLogSequenceNumber(2))
	})
}

func TestLog_EmptyReader(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_EmptyReader)
}

func testLog_EmptyReader(t *testing.T) {
	l := createLog(t, 320, 2)

	for _, err := range l.Iterate(record.NewLogSequenceNumber(0)) {
		if err != segment.ErrNoSegmentFile {
			t.Error("should have no files")
		}
	}
}

func TestLog_Iterate(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Iterate)
}

func testLog_Iterate(t *testing.T) {
	ctx := t.Context()

	l := createLogStart(t, 320, 2)
	tid := value.NewTenantID()

	records, _ := createRecords(tid, record.NewLogSequenceNumber(100))
	for _, rec := range records {
		if _, err := l.Update(ctx,
			tid,
			rec.TransactionID(),
			rec.Time().Value(),
			rec.Data(),
			rec.Reverse()); err != nil {
			t.Fatal(err)
		}
	}

	if err := l.Flush(ctx); err != nil {
		t.Fatal(err)
	}

	i := 0
	for rc, err := range l.Iterate(record.NewLogSequenceNumber(0)) {
		if err != nil {
			t.Fatal(err)
		}

		rec := records[i]

		if !reflect.DeepEqual(rc, rec) {
			t.Errorf("records do not match: %v, expected: %v",
				rc.LogSequenceNumber(),
				rec.LogSequenceNumber())
		}
		i++
	}
	if i != 100 {
		t.Errorf("incorrect count: %v, expected: %v", i, 100)
	}
}

func TestLog_Recover(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Recover)
}

func testLog_Recover(t *testing.T) {
	ctx := t.Context()

	fsys, dir := createFiles(t)
	tid := value.NewTenantID()
	records, _ := createRecords(tid, record.NewLogSequenceNumber(2))
	r0 := records[0]
	r1 := records[1]

	{ // "should append and flush"
		l := NewLog(256, 2, fsys, dir)
		if err := l.Open(t.Context()); err != nil {
			t.Fatal(err)
		}

		if err := l.StartWriter(); err != nil {
			t.Fatal(err)
		}

		if _, err := l.Update(ctx,
			tid,
			r0.TransactionID(),
			r0.Time().Value(),
			r0.Data(),
			r0.Reverse()); err != nil {
			t.Fatal(err)
		}

		if err := l.Flush(ctx); err != nil {
			t.Fatal(err)
		}

		if _, err := l.Update(ctx,
			tid,
			r1.TransactionID(),
			r1.Time().Value(),
			r1.Data(),
			r1.Reverse()); err != nil {
			t.Fatal(err)
		}

		if err := l.Flush(ctx); err != nil {
			t.Fatal(err)
		}

		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	}

	{ //"should recover"
		l := NewLog(256, 2, fsys, dir)
		if err := l.Open(t.Context()); err != nil {
			t.Fatal(err)
		}

		it, err := l.Recover()
		if err != nil {
			t.Fatal(err)
		}

		i := 0
		for rc := range it {
			if err != nil {
				t.Fatal(err)
			}

			rec := records[i]

			if !reflect.DeepEqual(rc, rec) {
				t.Errorf("records do not match: %v, expected: %v",
					rc.LogSequenceNumber(),
					rec.LogSequenceNumber())
			}
			i++
		}
		if i != 2 {
			t.Errorf("incorrect count: %v, expected: %v", i, 2)
		}
	}
}

func TestLog_RecoverMany(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_RecoverMany)
}

func testLog_RecoverMany(t *testing.T) {
	ctx := t.Context()

	fsys, dir := createFiles(t)
	tid := value.NewTenantID()

	var aCount1 = record.NewLogSequenceNumber(1)
	var aCount2 = record.NewLogSequenceNumber(1)
	aCount := aCount1.Value() + aCount2.Value()
	records1, _ := createRecords(tid, aCount1)
	records2, _ := createRecords(tid, aCount2)
	records := append(records1, records2...)

	{ // "should append and flush"
		l := NewLog(256, 2, fsys, dir)
		if err := l.Open(t.Context()); err != nil {
			t.Fatal(err)
		}

		if err := l.StartWriter(); err != nil {
			t.Fatal(err)
		}

		for _, rec := range records1 {
			if _, err := l.Update(ctx,
				tid,
				rec.TransactionID(),
				rec.Time().Value(),
				rec.Data(),
				rec.Reverse()); err != nil {
				t.Fatal(err)
			}
		}

		if err := l.Flush(ctx); err != nil {
			t.Fatal(err)
		}

		i := 0
		for rc, err := range l.Iterate(record.NewLogSequenceNumber(0)) {
			if err != nil {
				t.Fatal(err)
			}

			rec := records1[i]

			if !reflect.DeepEqual(rc, rec) {
				t.Errorf("records do not match: %v, expected: %v",
					rc.LogSequenceNumber(),
					rec.LogSequenceNumber())
			}
			i++
		}
		if i != int(aCount1.Value()) {
			t.Errorf("incorrect count: %v, expected: %v", i, aCount1)
		}

		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	}

	{ // "should append and flush more and iterate"
		l := NewLog(256, 2, fsys, dir)
		if err := l.Open(t.Context()); err != nil {
			t.Fatal(err)
		}

		if err := l.StartWriter(); err != nil {
			t.Fatal(err)
		}

		for _, rec := range records2 {
			if _, err := l.Update(ctx,
				tid,
				rec.TransactionID(),
				rec.Time().Value(),
				rec.Data(),
				rec.Reverse()); err != nil {
				t.Fatal(err)
			}
		}

		if err := l.Flush(ctx); err != nil {
			t.Fatal(err)
		}

		i := 0
		for rc, err := range l.Iterate(record.NewLogSequenceNumber(0)) {
			if err != nil {
				t.Fatal(err)
			}

			rec := records[i]

			if !reflect.DeepEqual(rc, rec) {
				t.Errorf("records do not match: %v, expected: %v",
					rc.LogSequenceNumber(),
					rec.LogSequenceNumber())
			}
			i++
		}
		if i != int(aCount) {
			t.Errorf("incorrect count: %v, expected: %v", i, aCount)
		}

		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLog_Reverse(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Reverse)
}

func testLog_Reverse(t *testing.T) {
	ctx := t.Context()

	l := createLogStart(t, 320, 2)
	tid := value.NewTenantID()

	records, _ := createRecords(tid, record.NewLogSequenceNumber(100))
	for _, rec := range records {
		if _, err := l.Update(ctx,
			tid,
			rec.TransactionID(),
			rec.Time().Value(),
			rec.Data(),
			rec.Reverse()); err != nil {
			t.Fatal(err)
		}
	}

	if err := l.Flush(ctx); err != nil {
		t.Fatal(err)
	}

	slices.Reverse(records)
	i := 0
	for rc, err := range l.Reverse() {
		if err != nil {
			t.Fatal(err)
		}

		rec := records[i]

		if !reflect.DeepEqual(rc, rec) {
			t.Errorf("records do not match: %v, expected: %v",
				rc.LogSequenceNumber(),
				rec.LogSequenceNumber())
		}
		i++
	}
	if i != 100 {
		t.Errorf("incorrect count: %v, expected: %v", i, 100)
	}
}

func createLogStart(t *testing.T, pageSize int64, segmentSize action.PageID) *Log {
	t.Helper()

	l := createLog(t, pageSize, segmentSize)
	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	return l
}

func createLog(t *testing.T, pageSize int64, segmentSize action.PageID) *Log {
	t.Helper()

	fsys, dir := createFiles(t)
	l := NewLog(pageSize, segmentSize, fsys, dir)
	if err := l.Open(t.Context()); err != nil {
		t.Fatal(err)
	}

	cleanupLog(t, l)

	return l
}

func cleanupLog(t *testing.T, l *Log) {
	t.Helper()

	t.Cleanup(func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	})
}

func createFiles(t *testing.T) (file.FileSystem, string) {
	// return &disk.FileSystem{}, t.TempDir()
	return filetesting.NewMockFileSystem(t, nil), "."
}

func createRecords(tid value.TenantID, count record.LogSequenceNumber) ([]*record.Record, record.LogSequenceNumber) {
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	records := make([]*record.Record, 0, count.Value())
	for i := range count.Value() {
		records = append(records, record.New(
			record.NewLogSequenceNumber(i+1),
			tid,
			record.NewTransactionID(2),
			record.NewTime(time.UnixMicro(1234567890)),
			record.ActionUpdate,
			data,
			reverse))
	}
	return records, count.Decrement()
}

func testPosition(t *testing.T, l *Log, lw, hw record.LogSequenceNumber) {
	t.Helper()

	if h := l.HighWater(); h != hw {
		t.Errorf("incorrect high water: %v, expected: %v", h, hw)
	}

	if l := l.LowWater(); l != lw {
		t.Errorf("incorrect low water: %v, expected: %v", l, lw)
	}
}

func runTest(t *testing.T, message string, f func(t *testing.T)) bool {
	return t.Run(message, func(t *testing.T) { t.Parallel(); synctest.Test(t, f) })
}
