package log

import (
	"reflect"
	"slices"
	"testing"
	"testing/fstest"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

func TestLog_Default(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Default)
}

func testLog_Default(t *testing.T) {
	wr := createLog(t)
	defer wr.Close()

	testPosition(t, wr, 0, 0)
}

func TestLog_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Append)
}

func testLog_Append(t *testing.T) {
	ctx := t.Context()

	wr := createLog(t)
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	if lsn, err := wr.Append(ctx, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse); err != nil {
		t.Error(err)
	} else if lsn != 1 {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, wr, 0, 1)
}

func TestLog_Flush(t *testing.T) {
	t.Parallel()

	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	runTest(t, "should flush", func(t *testing.T) {
		ctx := t.Context()

		wr := createLog(t)

		lsn1, err := wr.Append(ctx, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)
		if err != nil {
			t.Error(err)
		}

		_, err = wr.Append(ctx, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)
		if err != nil {
			t.Error(err)
		}

		if err := wr.Flush(lsn1); err != nil {
			t.Error(err)
		}

		testPosition(t, wr, 1, 2)
	})

	runTest(t, "should not flush beyond HighWater", func(t *testing.T) {
		ctx := t.Context()

		wr := createLog(t)

		_, err := wr.Append(ctx, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)
		if err != nil {
			t.Error(err)
		}

		_, err = wr.Append(ctx, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)
		if err != nil {
			t.Error(err)
		}

		if err := wr.Flush(10); err != nil {
			t.Error(err)
		}

		testPosition(t, wr, 2, 2)
	})

	runTest(t, "should write to multiple pages", func(t *testing.T) {
		ctx := t.Context()

		wr := createLog(t)

		count := 10

		for range count - 1 {
			_, err := wr.Append(ctx, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)
			if err != nil {
				t.Fatal(err)
			}
		}

		lsn2, err := wr.Append(ctx, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)
		if err != nil {
			t.Fatal(err)
		}

		if err := wr.Flush(lsn2); err != nil {
			t.Fatal(err)
		}

		if p := wr.PageID(); p != 2 {
			t.Errorf("incorrect value: %v, expected: %v", p, 2)
		}
	})

	runTest(t, "should return error if appending beyond maximum", func(t *testing.T) {
		t.Skip()
		// TODO: Test this
	})

	runTest(t, "should write after flushing", func(t *testing.T) {
		ctx := t.Context()

		wr := createLog(t)

		lsn1, err := wr.Append(ctx, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)
		if err != nil {
			t.Error(err)
		}

		if err := wr.Flush(lsn1); err != nil {
			t.Error(err)
		}

		lsn2, err := wr.Append(ctx, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse)
		if err != nil {
			t.Error(err)
		}

		if err := wr.Flush(lsn2); err != nil {
			t.Error(err)
		}

		testPosition(t, wr, 2, 2)
	})
}

func TestLog_EmptyReader(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_EmptyReader)
}

func testLog_EmptyReader(t *testing.T) {
	fsys, dir := createFiles(t)

	l := NewLog(256, 2, fsys, dir)
	if err := l.Open(t.Context()); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}()

	for _, err := range l.Iterate(0) {
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

	fsys, dir := createFiles(t)

	l := NewLog(256, 2, fsys, dir)
	if err := l.Open(t.Context()); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}()

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	records, _ := createRecords(100)
	var lsn record.LogSequenceNumber
	var err error
	for _, rec := range records {
		lsn, err = l.Append(ctx, rec.TransactionID(), rec.Time(), rec.Action(), rec.Data(), rec.Reverse())
		if err != nil {
			t.Fatal(err)
		}
	}

	if err = l.Flush(lsn); err != nil {
		t.Fatal(err)
	}

	i := 0
	for rc, err := range l.Iterate(0) {
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
	ctx := t.Context()

	fsys, dir := createFiles(t)
	records, _ := createRecords(2)
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

		lsn1, err := l.Append(ctx, r0.TransactionID(), r0.Time(), r0.Action(), r0.Data(), r0.Reverse())
		if err != nil {
			t.Fatal(err)
		}

		if err := l.Flush(lsn1); err != nil {
			t.Fatal(err)
		}

		lsn2, err := l.Append(ctx, r1.TransactionID(), r1.Time(), r1.Action(), r1.Data(), r1.Reverse())
		if err != nil {
			t.Fatal(err)
		}

		if err := l.Flush(lsn2); err != nil {
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
	ctx := t.Context()

	fsys, dir := createFiles(t)

	var aCount1 record.LogSequenceNumber = 1
	var aCount2 record.LogSequenceNumber = 1
	aCount := aCount1 + aCount2
	records1, _ := createRecords(aCount1)
	records2, _ := createRecords(aCount2)
	records := append(records1, records2...)

	{ // "should append and flush"
		l := NewLog(256, 2, fsys, dir)
		if err := l.Open(t.Context()); err != nil {
			t.Fatal(err)
		}

		if err := l.StartWriter(); err != nil {
			t.Fatal(err)
		}

		var lsn record.LogSequenceNumber
		var err error
		for _, rec := range records1 {
			lsn, err = l.Append(ctx, rec.TransactionID(), rec.Time(), rec.Action(), rec.Data(), rec.Reverse())
			if err != nil {
				t.Fatal(err)
			}
		}

		if err = l.Flush(lsn); err != nil {
			t.Fatal(err)
		}

		i := 0
		for rc, err := range l.Iterate(0) {
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
		if i != int(aCount1) {
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

		var lsn record.LogSequenceNumber
		var err error
		for _, rec := range records2 {
			lsn, err = l.Append(ctx, rec.TransactionID(), rec.Time(), rec.Action(), rec.Data(), rec.Reverse())
			if err != nil {
				t.Fatal(err)
			}
		}

		if err := l.Flush(lsn); err != nil {
			t.Fatal(err)
		}

		i := 0
		for rc, err := range l.Iterate(0) {
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

	fsys, dir := createFiles(t)

	l := NewLog(256, 2, fsys, dir)
	if err := l.Open(t.Context()); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}()

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	records, _ := createRecords(100)
	var lsn record.LogSequenceNumber
	var err error
	for _, rec := range records {
		lsn, err = l.Append(ctx, rec.TransactionID(), rec.Time(), rec.Action(), rec.Data(), rec.Reverse())
		if err != nil {
			t.Fatal(err)
		}
	}

	if err = l.Flush(lsn); err != nil {
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

func createLog(t *testing.T) *Log {
	t.Helper()

	fsys, dir := createFiles(t)
	l := NewLog(256, 3, fsys, dir)
	if err := l.Open(t.Context()); err != nil {
		t.Fatal(err)
	}

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	return l
}

func createFiles(t *testing.T) (file.FileSystem, string) {
	// return &disk.FileSystem{}, t.TempDir()
	return &mock.FileSystem{MapFS: fstest.MapFS{}}, "."
}

func createRecords(count record.LogSequenceNumber) ([]*record.Record, record.LogSequenceNumber) {
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	records := make([]*record.Record, 0, count)
	for i := range count {
		records = append(records, record.New(i+1, 2, time.UnixMicro(1234567890), record.ActionInsert, data, reverse))
	}
	return records, count - 1
}

func testPosition(t *testing.T, l *Log, lw, hw record.LogSequenceNumber) {
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
