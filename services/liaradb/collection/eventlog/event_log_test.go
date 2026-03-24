package eventlog

import (
	"path"
	"reflect"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/google/uuid"
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestEventLog_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventLog_Append)
}

func testEventLog_Append(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 2, 1024)
	el := New(s, btree.NewCursor(s))
	tn := tablename.NewFromString(path.Join(t.TempDir(), "testfile"))
	pid := value.NewPartitionID(0)

	records := []*entity.Event{{
		GlobalVersion: value.NewGlobalVersion(0),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(0),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(1),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(1),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(2),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(2),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(3),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(3),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(4),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(4),
		Data:          value.NewData([]byte{}),
	}}

	for _, r := range records {
		if _, err := el.Append(ctx, tn, pid, r); err != nil {
			t.Fatal(err)
		}
	}

	result := make([]*entity.Event, 0)

	for e, err := range el.Events(ctx, tn, pid) {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, e)
	}

	if !slices.EqualFunc(result, records, func(a, b *entity.Event) bool {
		return reflect.DeepEqual(a, b)
	}) {
		t.Errorf("incorrect result: %v, expected: %v", result, records)
	}
}

func TestEventLog_Find(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventLog_Find)
}

func testEventLog_Find(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 2, 1024)
	el := New(s, btree.NewCursor(s))
	tn := tablename.NewFromString(path.Join(t.TempDir(), "testfile"))
	pid := value.NewPartitionID(0)

	records := []*entity.Event{{
		GlobalVersion: value.NewGlobalVersion(0),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(0),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(1),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(1),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(2),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(2),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(3),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(3),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(4),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(4),
		Data:          value.NewData([]byte{}),
	}}

	for _, r := range records {
		if _, err := el.Append(ctx, tn, pid, r); err != nil {
			t.Fatal(err)
		}
	}

	e, err := el.Find(ctx, tn, pid, records[2].ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(e, records[2]) {
		t.Errorf("incorrect event: %v, expected: %v", e, records[2])
	}

	synctest.Wait()
}

func TestEventLog_GetAggregate(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventLog_GetAggregate)
}

func testEventLog_GetAggregate(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 2, 1024)
	el := New(s, btree.NewCursor(s))
	tn := tablename.NewFromString(path.Join(t.TempDir(), "testfile"))

	aggregateID := value.NewAggregateID(uuid.NewString())

	records := []*entity.Event{{
		GlobalVersion: value.NewGlobalVersion(0),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(1),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(1),
		ID:            value.NewEventID(),
		AggregateID:   aggregateID,
		Version:       value.NewVersion(1),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(2),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(2),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(3),
		ID:            value.NewEventID(),
		AggregateID:   aggregateID,
		Version:       value.NewVersion(2),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(4),
		ID:            value.NewEventID(),
		Version:       value.NewVersion(3),
		Data:          value.NewData([]byte{}),
	}}

	pid := value.NewPartitionID(0)
	for _, r := range records {
		if _, err := el.Append(ctx, tn, pid, r); err != nil {
			t.Fatal(err)
		}
	}

	want := []entity.Event{*records[1], *records[3]}

	result := make([]entity.Event, 0)
	for e, err := range el.GetAggregate(ctx, tn, pid, aggregateID) {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, *e)
	}

	if !slices.EqualFunc(result, want, func(a, b entity.Event) bool {
		return reflect.DeepEqual(a, b)
	}) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}

	synctest.Wait()
}

func TestEventLog_AppendEvent(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventLog_AppendEvent)
}

func testEventLog_AppendEvent(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 2, 1024)
	el := New(s, btree.NewCursor(s))
	tn := tablename.NewFromString(path.Join(t.TempDir(), "testfile"))
	pid := value.NewPartitionID(0)

	records := [][]byte{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10}}

	for i, r := range records {
		k := key.NewKey2([]byte(""), int64(i))
		if _, err := el.AppendEvent(ctx, tn, pid, k, buffer.NewFromSlice(r)); err != nil {
			t.Fatal(err)
		}
	}

	result := make([][]byte, 0)

	for n, err := range el.Iterate(ctx, tn, pid) {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, n)
	}

	if !slices.EqualFunc(result, records, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, records)
	}

	synctest.Wait()
}
