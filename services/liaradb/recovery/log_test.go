package recovery

import (
	"reflect"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/filecache"
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
	l := createLogStart(t, 320, 3, 320)

	testPosition(t, l, record.NewLogSequenceNumber(0), record.NewLogSequenceNumber(0))
}

func TestLog_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Append)
}

func testLog_Append(t *testing.T) {
	ctx := t.Context()

	l := createLogStart(t, 320, 3, 320)
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	if lsn, err := l.Update(ctx,
		value.NewTenantID(),
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
		record.CollectionValue,
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

	l := createLogStart(t, 320, 3, 320)
	var data = make([]byte, 0, 1024)
	for i := range 1024 {
		data = append(data, byte(i%255))
	}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	if _, err := l.Update(ctx,
		value.NewTenantID(),
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
		record.CollectionValue,
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

	t.Run("should flush", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			l := createLogStart(t, 320, 3, 320)
			tid := value.NewTenantID()

			if _, err := l.Update(ctx,
				tid,
				record.NewTransactionID(2),
				time.UnixMicro(1234567890),
				record.CollectionValue,
				data,
				reverse,
			); err != nil {
				t.Error(err)
			}

			testPosition(t, l, record.NewLogSequenceNumber(0), record.NewLogSequenceNumber(1))

			if _, err := l.Update(ctx,
				tid,
				record.NewTransactionID(2),
				time.UnixMicro(1234567890),
				record.CollectionValue,
				data,
				reverse,
			); err != nil {
				t.Error(err)
			}

			testPosition(t, l, record.NewLogSequenceNumber(0), record.NewLogSequenceNumber(2))

			time.Sleep(1 * time.Second)

			testPosition(t, l, record.NewLogSequenceNumber(2), record.NewLogSequenceNumber(2))
		})
	})

	t.Run("should not flush beyond HighWater", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			l := createLogStart(t, 320, 3, 320)
			tid := value.NewTenantID()

			if _, err := l.Update(ctx,
				tid,
				record.NewTransactionID(2),
				time.UnixMicro(1234567890),
				record.CollectionValue,
				data,
				reverse,
			); err != nil {
				t.Error(err)
			}

			if _, err := l.Update(ctx,
				tid,
				record.NewTransactionID(2),
				time.UnixMicro(1234567890),
				record.CollectionValue,
				data,
				reverse,
			); err != nil {
				t.Error(err)
			}

			time.Sleep(1 * time.Second)

			testPosition(t, l, record.NewLogSequenceNumber(2), record.NewLogSequenceNumber(2))
		})
	})

	t.Run("should write to multiple pages", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			l := createLogStart(t, 352, 4, 352)
			tid := value.NewTenantID()
			count := 14
			for range count {
				if _, err := l.Update(ctx,
					tid,
					record.NewTransactionID(2),
					time.UnixMicro(1234567890),
					record.CollectionValue,
					data,
					reverse,
				); err != nil {
					t.Fatal(err)
				}
			}

			time.Sleep(1 * time.Second)

			if p := l.PageID(); p != 3 {
				t.Errorf("incorrect value: %v, expected: %v", p, 3)
			}
		})
	})

	t.Run("should return error if appending beyond maximum", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			l := createLogStart(t, 32, 1, 32)

			if _, err := l.Update(ctx,
				value.NewTenantID(),
				record.NewTransactionID(2),
				time.UnixMicro(1234567890),
				record.CollectionValue,
				data,
				reverse,
			); err != raw.ErrInsufficientSpace {
				t.Fatal("should return error")
			}
		})
	})

	t.Run("should write after flushing", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			l := createLogStart(t, 320, 3, 320)
			tid := value.NewTenantID()

			if _, err := l.Update(ctx,
				tid,
				record.NewTransactionID(2),
				time.UnixMicro(1234567890),
				record.CollectionValue,
				data,
				reverse,
			); err != nil {
				t.Error(err)
			}

			time.Sleep(1 * time.Second)

			if _, err := l.Update(ctx,
				tid,
				record.NewTransactionID(2),
				time.UnixMicro(1234567890),
				record.CollectionValue,
				data,
				reverse,
			); err != nil {
				t.Error(err)
			}

			time.Sleep(1 * time.Second)

			testPosition(t, l, record.NewLogSequenceNumber(2), record.NewLogSequenceNumber(2))
		})
	})
}

func TestLog_FlushCheckpoint(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_FlushCheckpoint)
}

func testLog_FlushCheckpoint(t *testing.T) {
	ctx := t.Context()
	fsys, dir := createFiles()
	l := createLogAllStart(t, 320, 3, 320, fsys, dir)
	tid := value.NewTenantID()

	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	if _, err := l.Update(ctx,
		tid,
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
		record.CollectionValue,
		data,
		reverse,
	); err != nil {
		t.Fatal(err)
	}

	now := time.UnixMicro(1234567891)
	txid := record.NewTransactionID(1)

	_, err := l.FlushCheckpoint(now, txid)
	if err != nil {
		t.Fatal(err)
	}

	l1 := createLogAllStart(t, 320, 3, 320, fsys, dir)
	it, err := l1.Recover()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for range it {
		count++
	}

	if count != 0 {
		t.Errorf("incorrect count: %v, expected: %v", count, 2)
	}
}

func TestLog_EmptyReader(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_EmptyReader)
}

func testLog_EmptyReader(t *testing.T) {
	l := createLogStart(t, 320, 2, 320)

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

	l := createLogStart(t, 320, 2, 320)
	tid := value.NewTenantID()

	records, _ := createRecords(tid, record.NewLogSequenceNumber(100))
	for _, rec := range records {
		if _, err := l.Update(ctx,
			tid,
			rec.TransactionID(),
			rec.Time().Value(),
			rec.Collection(),
			rec.Data(),
			rec.Reverse(),
		); err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(1 * time.Second)

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

	fsys, dir := createFiles()
	tid := value.NewTenantID()
	records, _ := createRecords(tid, record.NewLogSequenceNumber(2))
	r0 := records[0]
	r1 := records[1]

	{ // "should append and flush"
		l := NewLog(256, 2, 256, fsys, dir)
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
			r0.Collection(),
			r0.Data(),
			r0.Reverse(),
		); err != nil {
			t.Fatal(err)
		}

		time.Sleep(1 * time.Second)

		if _, err := l.Update(ctx,
			tid,
			r1.TransactionID(),
			r1.Time().Value(),
			r1.Collection(),
			r1.Data(),
			r1.Reverse(),
		); err != nil {
			t.Fatal(err)
		}

		time.Sleep(1 * time.Second)

		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	}

	{ //"should recover"
		l := NewLog(256, 2, 256, fsys, dir)
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

	fsys, dir := createFiles()
	tid := value.NewTenantID()

	var aCount1 = record.NewLogSequenceNumber(1)
	var aCount2 = record.NewLogSequenceNumber(1)
	aCount := aCount1.Value() + aCount2.Value()
	records1, _ := createRecords(tid, aCount1)
	records2, _ := createRecords(tid, aCount2)
	records := append(records1, records2...)

	{ // "should append and flush"
		l := NewLog(256, 2, 256, fsys, dir)
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
				rec.Collection(),
				rec.Data(),
				rec.Reverse(),
			); err != nil {
				t.Fatal(err)
			}
		}

		time.Sleep(1 * time.Second)

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
		l := NewLog(256, 2, 256, fsys, dir)
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
				rec.Collection(),
				rec.Data(),
				rec.Reverse(),
			); err != nil {
				t.Fatal(err)
			}
		}

		time.Sleep(1 * time.Second)

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

	l := createLogStart(t, 320, 2, 320)
	tid := value.NewTenantID()

	records, _ := createRecords(tid, record.NewLogSequenceNumber(100))
	for _, rec := range records {
		if _, err := l.Update(ctx,
			tid,
			rec.TransactionID(),
			rec.Time().Value(),
			rec.Collection(),
			rec.Data(),
			rec.Reverse(),
		); err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(1 * time.Second)

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

func TestLog_Commit(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Commit)
}

func testLog_Commit(t *testing.T) {
	ctx := t.Context()

	fsys, dir := createFiles()
	l := createLogAllStart(t, 320, 3, 320, fsys, dir)

	if lsn, err := l.Commit(ctx,
		value.NewTenantID(),
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
	); err != nil {
		t.Error(err)
	} else if lsn != record.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, record.NewLogSequenceNumber(1), record.NewLogSequenceNumber(1))

	l2 := createLogAllStart(t, 320, 3, 320, fsys, dir)

	it, err := l2.Recover()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for r := range it {
		count++
		if a := r.Action(); a != record.ActionCommit {
			t.Errorf("incorrect action: %v, expected: %v", a, record.ActionCommit)
		}
	}

	if count != 1 {
		t.Errorf("incorrect count: %v, expected: %v", count, 1)
	}
}

func TestLog_Insert(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Insert)
}

func testLog_Insert(t *testing.T) {
	ctx := t.Context()

	fsys, dir := createFiles()
	l := createLogAllStart(t, 320, 3, 320, fsys, dir)
	var data = []byte{0, 1, 2, 3, 4, 5}

	if lsn, err := l.Insert(ctx,
		value.NewTenantID(),
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
		record.CollectionEvent,
		data,
	); err != nil {
		t.Error(err)
	} else if lsn != record.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, record.NewLogSequenceNumber(0), record.NewLogSequenceNumber(1))

	time.Sleep(1 * time.Second)

	l2 := createLogAllStart(t, 320, 3, 320, fsys, dir)

	it, err := l2.Recover()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for r := range it {
		count++
		if a := r.Action(); a != record.ActionInsert {
			t.Errorf("incorrect action: %v, expected: %v", a, record.ActionInsert)
		}
	}

	if count != 1 {
		t.Errorf("incorrect count: %v, expected: %v", count, 1)
	}
}

func TestLog_Rollback(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Rollback)
}

func testLog_Rollback(t *testing.T) {
	ctx := t.Context()

	fsys, dir := createFiles()
	l := createLogAllStart(t, 320, 3, 320, fsys, dir)

	if lsn, err := l.Rollback(ctx,
		value.NewTenantID(),
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
	); err != nil {
		t.Error(err)
	} else if lsn != record.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, record.NewLogSequenceNumber(1), record.NewLogSequenceNumber(1))

	l2 := createLogAllStart(t, 320, 3, 320, fsys, dir)

	it, err := l2.Recover()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for r := range it {
		count++
		if a := r.Action(); a != record.ActionRollback {
			t.Errorf("incorrect action: %v, expected: %v", a, record.ActionRollback)
		}
	}

	if count != 1 {
		t.Errorf("incorrect count: %v, expected: %v", count, 1)
	}
}

func TestLog_Start(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Start)
}

func testLog_Start(t *testing.T) {
	ctx := t.Context()

	fsys, dir := createFiles()
	l := createLogAllStart(t, 320, 3, 320, fsys, dir)

	if lsn, err := l.Start(ctx,
		value.NewTenantID(),
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
	); err != nil {
		t.Error(err)
	} else if lsn != record.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, record.NewLogSequenceNumber(0), record.NewLogSequenceNumber(1))

	time.Sleep(1 * time.Second)

	l2 := createLogAllStart(t, 320, 3, 320, fsys, dir)

	it, err := l2.Recover()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for r := range it {
		count++
		if a := r.Action(); a != record.ActionStart {
			t.Errorf("incorrect action: %v, expected: %v", a, record.ActionStart)
		}
	}

	if count != 1 {
		t.Errorf("incorrect count: %v, expected: %v", count, 1)
	}
}

func TestLog_Update(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Update)
}

func testLog_Update(t *testing.T) {
	ctx := t.Context()

	fsys, dir := createFiles()
	l := createLogAllStart(t, 320, 3, 320, fsys, dir)
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	if lsn, err := l.Update(ctx,
		value.NewTenantID(),
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
		record.CollectionValue,
		data,
		reverse,
	); err != nil {
		t.Error(err)
	} else if lsn != record.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, record.NewLogSequenceNumber(0), record.NewLogSequenceNumber(1))

	time.Sleep(1 * time.Second)

	l2 := createLogAllStart(t, 320, 3, 320, fsys, dir)

	it, err := l2.Recover()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for r := range it {
		count++
		if a := r.Action(); a != record.ActionUpdate {
			t.Errorf("incorrect action: %v, expected: %v", a, record.ActionUpdate)
		}
	}

	if count != 1 {
		t.Errorf("incorrect count: %v, expected: %v", count, 1)
	}
}

func createLogStart(t *testing.T,
	pageSize int64,
	segmentSize action.PageID,
	recordSize int64,
) *Log {
	t.Helper()

	fsys, dir := createFiles()
	l := createLog(t, pageSize, segmentSize, recordSize, fsys, dir)
	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	return l
}

func createLogAllStart(t *testing.T,
	pageSize int64,
	segmentSize action.PageID,
	recordSize int64,
	fsys filecache.FileSystem,
	dir string,
) *Log {
	t.Helper()

	l := createLog(t, pageSize, segmentSize, recordSize, fsys, dir)
	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	return l
}

func createLog(t *testing.T,
	pageSize int64,
	segmentSize action.PageID,
	recordSize int64,
	fsys filecache.FileSystem,
	dir string,
) *Log {
	t.Helper()

	l := NewLog(pageSize, segmentSize, recordSize, fsys, dir)
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

func createFiles() (filecache.FileSystem, string) {
	// return &disk.FileSystem{}, t.TempDir()
	return filetesting.New(nil), "."
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
			record.CollectionValue,
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
