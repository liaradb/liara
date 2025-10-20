package transaction

import (
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/storage/eventlog"
)

func TestTransaction_Insert(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTransaction_Insert)
}

func testTransaction_Insert(t *testing.T) {
	m, l := createManager(t)
	ctx := t.Context()

	tx := m.Next()

	if err := tx.Insert(ctx, "a", time.UnixMicro(1234567890), nil); err != nil {
		t.Fatal(err)
	}

	if err := l.Flush(ctx, tx.LogSequenceNumber()); err != nil {
		t.Fatal(err)
	}

	c := 0
	for rc, err := range l.Iterate(0) {
		if err != nil {
			t.Fatal(err)
		}

		if lsn := tx.LogSequenceNumber(); lsn != rc.LogSequenceNumber() {
			t.Errorf("lsn does not match: %v, expected: %v", lsn, rc.LogSequenceNumber())
		}

		c++
	}

	if c != 1 {
		t.Errorf("incorrect record count: %v, expected: %v", c, 1)
	}
}

func TestTransaction_Commit(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTransaction_Commit)
}

func testTransaction_Commit(t *testing.T) {
	m, l := createManager(t)
	ctx := t.Context()

	tx := m.Next()

	records := [][]byte{{1, 2, 3, 4, 5}}

	if err := tx.Insert(ctx, "a", time.UnixMicro(1234567890), records[0]); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(ctx, time.UnixMicro(1234567890)); err != nil {
		t.Fatal(err)
	}

	lsns := []record.LogSequenceNumber{1, 2}
	actions := []record.Action{record.ActionInsert, record.ActionCommit}

	c := 0
	for rc, err := range l.Iterate(0) {
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

	for b, err := range eventlog.New(m.storage).Iterate(ctx, "filename") {
		if err != nil {
			t.Fatal(err)
		}

		for i, err := range b.Items() {
			if err != nil {
				t.Fatal(err)
			}

			result = append(result, i)
		}
	}

	if !slices.EqualFunc(result, records, slices.Equal) {
		t.Errorf("incorrect records do not match: %v, expected: %v", result, records)
	}
}
