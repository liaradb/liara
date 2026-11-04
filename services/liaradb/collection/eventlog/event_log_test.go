package eventlog

import (
	"reflect"
	"slices"
	"testing"
	"testing/synctest"
	"uuid"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/util/testing/logtesting"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestEventLog_Append(t *testing.T) {
	storagetesting.SyncTest(t, 2, 1024, testEventLog_Append)
}

func testEventLog_Append(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	el := New(s.Storage, btree.NewCursor(s.Storage))
	tf := tablename.NewTenantFolder(value.NewTenantID())
	pid := value.NewPartitionID(0)
	tn := tf.EventLog(pid)

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
		if err := el.Append(ctx, &logtesting.MockLog{}, tn, r); err != nil {
			t.Fatal(err)
		}
	}

	result := make([]*entity.Event, 0)

	for e, err := range el.Events(ctx, &logtesting.MockLog{}, tn) {
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

func TestEventLog_EventsAfterGlobalVersion(t *testing.T) {
	t.Skip()
	storagetesting.SyncTest(t, 2, 1024, testEventLog_EventsAfterGlobalVersion)
}

func testEventLog_EventsAfterGlobalVersion(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	el := New(s.Storage, btree.NewCursor(s.Storage))
	tf := tablename.NewTenantFolder(value.NewTenantID())
	pid := value.NewPartitionID(0)
	tn := tf.EventLog(pid)

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
		if err := el.Append(ctx, &logtesting.MockLog{}, tn, r); err != nil {
			t.Fatal(err)
		}
	}

	partialResults := make([]*entity.Event, 0)

	for n, err := range el.EventsAfterGlobalVersion(ctx, tn, value.NewGlobalVersion(3)) {
		if err != nil {
			t.Fatal(err)
		}

		partialResults = append(partialResults, n)
	}

	if !slices.EqualFunc(partialResults, records[:2], func(a, b *entity.Event) bool {
		return reflect.DeepEqual(a, b)
	}) {
		t.Errorf("incorrect result: %v, expected: %v", partialResults, records[:2])
	}
}

func TestEventLog_Find(t *testing.T) {
	storagetesting.SyncTest(t, 2, 1024, testEventLog_Find)
}

func testEventLog_Find(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	el := New(s.Storage, btree.NewCursor(s.Storage))
	tf := tablename.NewTenantFolder(value.NewTenantID())
	pid := value.NewPartitionID(0)
	tn := tf.EventLog(pid)

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
		if err := el.Append(ctx, &logtesting.MockLog{}, tn, r); err != nil {
			t.Fatal(err)
		}
	}

	e, err := el.Find(ctx, tn, records[2].ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(e, records[2]) {
		t.Errorf("incorrect event: %v, expected: %v", e, records[2])
	}

	synctest.Wait()
}

func TestEventLog_GetAggregate(t *testing.T) {
	storagetesting.SyncTest(t, 2, 1024, testEventLog_GetAggregate)
}

func testEventLog_GetAggregate(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	el := New(s.Storage, btree.NewCursor(s.Storage))
	tf := tablename.NewTenantFolder(value.NewTenantID())
	pid := value.NewPartitionID(0)
	tn := tf.EventLog(pid)

	aggregateID := value.NewAggregateID(uuid.NewV7().String())

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

	for _, r := range records {
		if err := el.Append(ctx, &logtesting.MockLog{}, tn, r); err != nil {
			t.Fatal(err)
		}
	}

	want := []entity.Event{*records[1], *records[3]}

	result := make([]entity.Event, 0)
	for e, err := range el.GetAggregate(ctx, tn, aggregateID) {
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
