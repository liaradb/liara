package transaction

import (
	"errors"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/recovery/record"
)

func TestTransaction_Insert(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTransaction_Insert)
}

func testTransaction_Insert(t *testing.T) {
	m, l := createManager(t)
	ctx := t.Context()

	tid := value.NewTenantID()
	tx, err := m.Next(ctx, tid)
	if err != nil {
		t.Fatal(err)
	}

	if err := tx.Insert(ctx, tablename.NewFromString("a"), time.UnixMicro(1234567890), &entity.Event{}, nil); err != nil {
		t.Fatal(err)
	}

	if err := l.Flush(ctx); err != nil {
		t.Fatal(err)
	}

	c := 0
	for rc, err := range l.Iterate(record.NewLogSequenceNumber(0)) {
		if err != nil {
			t.Fatal(err)
		}

		if lsn := l.HighWater(); lsn != rc.LogSequenceNumber() {
			t.Errorf("lsn does not match: %v, expected: %v", lsn, rc.LogSequenceNumber())
		}

		c++
	}

	if c != 1 {
		t.Errorf("incorrect record count: %v, expected: %v", c, 1)
	}

	synctest.Wait()
}

func TestTransaction_Insert__Unique(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTransaction_Insert__Unique)
}

func testTransaction_Insert__Unique(t *testing.T) {
	m, _ := createManager(t)
	ctx := t.Context()

	tid := value.NewTenantID()
	tn := tablename.New(tid)
	id := value.NewAggregateID("b")
	version := value.NewVersion(1)
	tm := time.UnixMicro(1234567890)

	tx, err := m.Next(ctx, tid)
	if err != nil {
		t.Fatal(err)
	}

	if err := Run(ctx, tx, tm, func() error {
		return tx.Insert(ctx, tn, tm, &entity.Event{
			AggregateID: id,
			Version:     version,
		}, nil)
	}); err != nil {
		t.Fatal(err)
	}

	tx, err = m.Next(ctx, tid)
	if err != nil {
		t.Fatal(err)
	}

	if err := Run(ctx, tx, tm, func() error {
		return tx.Insert(ctx, tn, time.UnixMicro(1234567890), &entity.Event{
			AggregateID: id,
			Version:     version,
		}, nil)
	}); !errors.Is(err, btree.ErrExists) {
		t.Fatalf("incorrect error: %v, expected: %v", err, btree.ErrExists)
	}

	synctest.Wait()
}

func TestTransaction_Insert__UniqueCurrent(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTransaction_Insert__UniqueCurrent)
}

func testTransaction_Insert__UniqueCurrent(t *testing.T) {
	m, _ := createManager(t)
	ctx := t.Context()

	tid := value.NewTenantID()
	tx, err := m.Next(ctx, tid)
	if err != nil {
		t.Fatal(err)
	}

	tn := tablename.NewFromString("a")
	id := value.NewAggregateID("b")
	version := value.NewVersion(1)
	tm := time.UnixMicro(1234567890)

	if err := tx.Insert(ctx, tn, tm, &entity.Event{
		AggregateID: id,
		Version:     version,
	}, nil); err != nil {
		t.Fatal(err)
	}

	if err := tx.Insert(ctx, tn, time.UnixMicro(1234567890), &entity.Event{
		AggregateID: id,
		Version:     version,
	}, nil); !errors.Is(err, btree.ErrExists) {
		t.Fatalf("incorrect error: %v, expected: %v", err, btree.ErrExists)
	}

	synctest.Wait()
}

func TestTransaction_Commit(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTransaction_Commit)
}

func testTransaction_Commit(t *testing.T) {
	m, l := createManager(t)
	ctx := t.Context()

	tid := value.NewTenantID()
	tx, err := m.Next(ctx, tid)
	if err != nil {
		t.Fatal(err)
	}

	type item struct {
		e    *entity.Event
		data []byte
	}
	items := []item{{
		e:    &entity.Event{},
		data: []byte{1},
	}, {
		e:    &entity.Event{},
		data: []byte{2},
	}, {
		e:    &entity.Event{},
		data: []byte{3},
	}, {
		e:    &entity.Event{},
		data: []byte{4},
	}, {
		e:    &entity.Event{},
		data: []byte{5},
	}}

	tn := tablename.New(tid)
	pid := value.NewPartitionID(0)

	if err := tx.Insert(ctx, tn, time.UnixMicro(1234567890), items[0].e, items[0].data); err != nil {
		t.Fatal(err)
	}

	if err := tx.commit(ctx, time.UnixMicro(1234567890)); err != nil {
		t.Fatal(err)
	}

	lsns := []record.LogSequenceNumber{record.NewLogSequenceNumber(1), record.NewLogSequenceNumber(2)}
	actions := []record.Action{record.ActionInsert, record.ActionCommit}

	c := 0
	for rc, err := range l.Iterate(record.NewLogSequenceNumber(0)) {
		if err != nil {
			t.Fatal(err)
		}

		if lsn := rc.LogSequenceNumber(); lsn != lsns[c] {
			t.Errorf("lsn does not match: %v, expected: %v", lsn, lsns[c])
		}

		if a := rc.Action(); a != actions[c] {
			t.Errorf("action does not match: %v, expected: %v", a, actions[c])
		}

		c++
	}

	if c != 2 {
		t.Errorf("incorrect record count: %v, expected: %v", c, 2)
	}

	result := [][]byte{}

	for n, err := range m.collections.EventLog.Iterate(ctx, tn, pid) {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, n)
	}

	records := [][]byte{items[0].data}
	if !slices.EqualFunc(result, records, slices.Equal) {
		t.Errorf("incorrect records do not match: %v, expected: %v", result, records)
	}

	synctest.Wait()
}

func TestTransaction_Rollback(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTransaction_Rollback)
}

func testTransaction_Rollback(t *testing.T) {
	m, l := createManager(t)
	ctx := t.Context()

	tid := value.NewTenantID()
	tx, err := m.Next(ctx, tid)
	if err != nil {
		t.Fatal(err)
	}

	records := [][]byte{{1, 2, 3, 4, 5}}

	tn := tablename.NewFromString("a")
	pid := value.NewPartitionID(0)

	if err := tx.Insert(ctx, tn, time.UnixMicro(1234567890), &entity.Event{}, records[0]); err != nil {
		t.Fatal(err)
	}

	if err := tx.rollback(ctx, time.UnixMicro(1234567890)); err != nil {
		t.Fatal(err)
	}

	lsns := []record.LogSequenceNumber{record.NewLogSequenceNumber(1), record.NewLogSequenceNumber(2)}
	actions := []record.Action{record.ActionInsert, record.ActionRollback}

	c := 0
	for rc, err := range l.Iterate(record.NewLogSequenceNumber(0)) {
		if err != nil {
			t.Fatal(err)
		}

		if lsn := rc.LogSequenceNumber(); lsn != lsns[c] {
			t.Errorf("lsn does not match: %v, expected: %v", lsn, lsns[c])
		}

		if a := rc.Action(); a != actions[c] {
			t.Errorf("action does not match: %v, expected: %v", a, actions[c])
		}

		c++
	}

	if c != 2 {
		t.Errorf("incorrect record count: %v, expected: %v", c, 2)
	}

	result := [][]byte{}

	for n, err := range m.collections.EventLog.Iterate(ctx, tn, pid) {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, n)
	}

	if length := len(result); length != 0 {
		t.Errorf("incorrect result length: %v, expected: %v", length, 0)
	}

	synctest.Wait()
}
