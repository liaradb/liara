package log

import (
	"reflect"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/transaction/record"
	"github.com/liaradb/liaradb/util/testing/filetesting"
)

func TestLog_Default(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Default)
}

func testLog_Default(t *testing.T) {
	l := createLogStart(t, 320, 3, 320)

	testPosition(t, l, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(0))
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
	} else if lsn != logpage.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(1))
}

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

	testPosition(t, l, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(0))
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

			testPosition(t, l, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(1))

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

			testPosition(t, l, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(2))

			if err := l.Flush(t.Context()); err != nil {
				t.Error(err)
			}

			testPosition(t, l, logpage.NewLogSequenceNumber(2), logpage.NewLogSequenceNumber(2))
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

			if err := l.Flush(t.Context()); err != nil {
				t.Error(err)
			}

			testPosition(t, l, logpage.NewLogSequenceNumber(2), logpage.NewLogSequenceNumber(2))
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

			if err := l.Flush(t.Context()); err != nil {
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
			l := createLogStart(t, 48, 1, 2)

			if _, err := l.Update(ctx,
				value.NewTenantID(),
				record.NewTransactionID(2),
				time.UnixMicro(1234567890),
				record.CollectionValue,
				data,
				reverse,
			); err == nil {
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

			if err := l.Flush(t.Context()); err != nil {
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

			if err := l.Flush(t.Context()); err != nil {
				t.Error(err)
			}

			testPosition(t, l, logpage.NewLogSequenceNumber(2), logpage.NewLogSequenceNumber(2))
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

	_, err := l.Checkpoint(t.Context(), now, txid)
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

	it, err := l.Recover()
	if err != nil {
		t.Fatal(err)
	}

	c := 0
	for range it {
		c++
	}

	if c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}
}

func TestLog_Iterate(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Iterate)
}

func testLog_Iterate(t *testing.T) {
	ctx := t.Context()

	fsys, dir := createFiles()

	l := createLogAllStart(t, 320, 2, 320, fsys, dir)
	tid := value.NewTenantID()

	count := 100
	records, _ := createRecords(tid, uint64(count), 0)
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

	if err := l.Flush(t.Context()); err != nil {
		t.Error(err)
	}

	l.Close()

	synctest.Wait()

	l = createLogAllStart(t, 320, 2, 320, fsys, dir)

	it, err := l.Recover()
	if err != nil {
		t.Fatal(err)
	}

	i := 0
	for rc := range it {
		rec := records[i]

		if !reflect.DeepEqual(rc, rec) {
			t.Errorf("records do not match: %v, expected: %v",
				rc.LogSequenceNumber(),
				rec.LogSequenceNumber())
		}
		i++
	}
	if i != count {
		t.Errorf("incorrect count: %v, expected: %v", i, count)
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
	records, _ := createRecords(tid, 2, 0)
	r0 := records[0]
	r1 := records[1]

	{ // "should append and flush"
		l := New(256, 2, 256, 100, fsys, dir)
		if err := l.Run(t.Context()); err != nil {
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

		if err := l.Flush(t.Context()); err != nil {
			t.Error(err)
		}

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

		if err := l.Flush(t.Context()); err != nil {
			t.Error(err)
		}

		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	}

	{ //"should recover"
		l := New(256, 2, 256, 100, fsys, dir)
		if err := l.Run(t.Context()); err != nil {
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

	var aCount1 = logpage.NewLogSequenceNumber(1)
	var aCount2 = logpage.NewLogSequenceNumber(1)
	aCount := aCount1.Value() + aCount2.Value()
	records1, _ := createRecords(tid, 1, 0)
	records2, _ := createRecords(tid, 1, 1)
	records := append(records1, records2...)

	{ // "should append and flush"
		l := New(256, 2, 256, 100, fsys, dir)
		if err := l.Run(t.Context()); err != nil {
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

		if err := l.Flush(t.Context()); err != nil {
			t.Error(err)
		}

		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	}

	{ // "should iterate"
		l := New(256, 2, 256, 100, fsys, dir)
		if err := l.Run(t.Context()); err != nil {
			t.Fatal(err)
		}

		it, err := l.Recover()
		if err != nil {
			t.Fatal(err)
		}

		i := 0
		for rc := range it {
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

		synctest.Wait()
	}

	{ // "should append and flush more"
		l := New(256, 2, 256, 100, fsys, dir)
		if err := l.Run(t.Context()); err != nil {
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

		if err := l.Flush(t.Context()); err != nil {
			t.Error(err)
		}

		if err := l.Close(); err != nil {
			t.Fatal(err)
		}

		synctest.Wait()
	}

	{ // "should iterate"
		l := New(256, 2, 256, 100, fsys, dir)
		if err := l.Run(t.Context()); err != nil {
			t.Fatal(err)
		}

		it, err := l.Recover()
		if err != nil {
			t.Fatal(err)
		}

		i := 0
		for rc := range it {

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

func TestLog_Recover__Iterate(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_Recover__Iterate)
}

func testLog_Recover__Iterate(t *testing.T) {
	ctx := t.Context()

	fsys, dir := createFiles()

	l := createLogAllStart(t, 320, 2, 320, fsys, dir)
	tid := value.NewTenantID()

	count := 100
	records, _ := createRecords(tid, uint64(count), 0)
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

	if err := l.Flush(t.Context()); err != nil {
		t.Error(err)
	}

	l.Close()

	synctest.Wait()

	l = createLogAllStart(t, 320, 2, 320, fsys, dir)

	it, err := l.Recover()
	if err != nil {
		t.Fatal(err)
	}

	i := 0
	for rc := range it {
		rec := records[i]

		if !reflect.DeepEqual(rc, rec) {
			t.Errorf("records do not match: %v, expected: %v",
				rc.LogSequenceNumber(),
				rec.LogSequenceNumber())
		}
		i++
	}
	if i != count {
		t.Errorf("incorrect count: %v, expected: %v", i, count)
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
	tid := value.NewTenantID()
	txid1 := record.NewTransactionID(1)
	txid2 := record.NewTransactionID(2)

	if _, err := l.Start(ctx,
		tid,
		txid1,
		time.UnixMicro(1234567890)); err != nil {
		t.Fatal(err)
	}

	if _, err := l.Insert(ctx,
		tid,
		txid1,
		time.UnixMicro(1234567890),
		record.CollectionSystem,
		make([]byte, 200)); err != nil {
		t.Fatal(err)
	}

	if lsn, err := l.Commit(ctx,
		tid,
		txid1,
		time.UnixMicro(1234567890),
	); err != nil {
		t.Error(err)
	} else if lsn != logpage.NewLogSequenceNumber(3) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 3)
	}

	if _, err := l.Start(ctx,
		tid,
		txid2,
		time.UnixMicro(1234567890)); err != nil {
		t.Fatal(err)
	}

	if _, err := l.Insert(ctx,
		tid,
		txid2,
		time.UnixMicro(1234567890),
		record.CollectionSystem,
		make([]byte, 200)); err != nil {
		t.Fatal(err)
	}

	if lsn, err := l.Commit(ctx,
		tid,
		txid2,
		time.UnixMicro(1234567890),
	); err != nil {
		t.Error(err)
	} else if lsn != logpage.NewLogSequenceNumber(6) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 6)
	}

	synctest.Wait()

	testPosition(t, l, logpage.NewLogSequenceNumber(6), logpage.NewLogSequenceNumber(6))

	l2 := createLogAllStart(t, 320, 3, 320, fsys, dir)

	it, err := l2.Recover()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for r := range it {
		switch count {
		case 0:
			if a := r.Action(); a != record.ActionStart {
				t.Errorf("incorrect action: %v, expected: %v", a, record.ActionStart)
			}
		case 1:
			if a := r.Action(); a != record.ActionInsert {
				t.Errorf("incorrect action: %v, expected: %v", a, record.ActionInsert)
			}
		case 2:
			if a := r.Action(); a != record.ActionCommit {
				t.Errorf("incorrect action: %v, expected: %v", a, record.ActionCommit)
			}
		case 3:
			if a := r.Action(); a != record.ActionStart {
				t.Errorf("incorrect action: %v, expected: %v", a, record.ActionStart)
			}
		case 4:
			if a := r.Action(); a != record.ActionInsert {
				t.Errorf("incorrect action: %v, expected: %v", a, record.ActionInsert)
			}
		case 5:
			if a := r.Action(); a != record.ActionCommit {
				t.Errorf("incorrect action: %v, expected: %v", a, record.ActionCommit)
			}
		}
		count++
	}

	if count != 6 {
		t.Errorf("incorrect count: %v, expected: %v", count, 6)
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
	} else if lsn != logpage.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(1))

	if err := l.Flush(t.Context()); err != nil {
		t.Error(err)
	}

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

func TestLog_InsertAndCommit(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testLog_InsertAndCommit)
}

func testLog_InsertAndCommit(t *testing.T) {
	ctx := t.Context()

	fsys, dir := createFiles()
	l := createLogAllStart(t, 320, 3, 320, fsys, dir)
	var data = []byte{0, 1, 2, 3, 4, 5}

	wg := sync.WaitGroup{}
	wg.Go(func() {
		time.Sleep(1 * time.Second)
		if err := l.Flush(t.Context()); err != nil {
			t.Error(err)
		}
	})

	if lsn, err := l.Insert(ctx,
		value.NewTenantID(),
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
		record.CollectionEvent,
		data,
	); err != nil {
		t.Error(err)
	} else if lsn != logpage.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(1))

	if lsn, err := l.Commit(ctx,
		value.NewTenantID(),
		record.NewTransactionID(2),
		time.UnixMicro(1234567890),
	); err != nil {
		t.Error(err)
	} else if lsn != logpage.NewLogSequenceNumber(2) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 2)
	}

	wg.Wait()

	l2 := createLogAllStart(t, 320, 3, 320, fsys, dir)

	it, err := l2.Recover()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for r := range it {
		switch count {
		case 0:
			if a := r.Action(); a != record.ActionInsert {
				t.Errorf("incorrect action: %v, expected: %v", a, record.ActionInsert)
			}
		case 1:
			if a := r.Action(); a != record.ActionCommit {
				t.Errorf("incorrect action: %v, expected: %v", a, record.ActionCommit)
			}
		}
		count++
	}

	if count != 2 {
		t.Errorf("incorrect count: %v, expected: %v", count, 2)
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
	} else if lsn != logpage.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	synctest.Wait()

	testPosition(t, l, logpage.NewLogSequenceNumber(1), logpage.NewLogSequenceNumber(1))

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
	} else if lsn != logpage.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(1))

	if err := l.Flush(t.Context()); err != nil {
		t.Error(err)
	}

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
	} else if lsn != logpage.NewLogSequenceNumber(1) {
		t.Errorf("incorrect value: %v, expected: %v", lsn, 1)
	}

	testPosition(t, l, logpage.NewLogSequenceNumber(0), logpage.NewLogSequenceNumber(1))

	if err := l.Flush(t.Context()); err != nil {
		t.Error(err)
	}

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
	segmentSize segment.PageID,
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
	segmentSize segment.PageID,
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
	segmentSize segment.PageID,
	recordSize int64,
	fsys filecache.FileSystem,
	dir string,
) *Log {
	t.Helper()

	l := New(pageSize, segmentSize, recordSize, 100, fsys, dir)
	if err := l.Run(t.Context()); err != nil {
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

func createRecords(tid value.TenantID, count, offset uint64) ([]*record.Record, logpage.LogSequenceNumber) {
	var data = []byte{0, 1, 2, 3, 4, 5}
	var reverse = []byte{6, 7, 8, 9, 10, 11}

	records := make([]*record.Record, 0, count)
	for i := range count {
		rc := record.New(
			tid,
			record.NewTransactionID(2),
			record.NewTime(time.UnixMicro(1234567890)),
			record.ActionUpdate,
			record.CollectionValue,
			data,
			reverse)
		rc.SetLogSequenceNumber(logpage.NewLogSequenceNumber(i + 1 + offset))
		records = append(records, rc)
	}
	return records, logpage.NewLogSequenceNumber(count).Decrement()
}

func testPosition(t *testing.T, l *Log, lw, hw logpage.LogSequenceNumber) {
	t.Helper()

	if h := l.HighWater(); h != hw {
		t.Errorf("incorrect high water: %v, expected: %v", h, hw)
	}

	if l := l.LowWater(); l != lw {
		t.Errorf("incorrect low water: %v, expected: %v", l, lw)
	}
}
