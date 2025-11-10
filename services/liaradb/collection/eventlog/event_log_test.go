package eventlog

import (
	"path"
	"reflect"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/google/uuid"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
)

func TestEventLog_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventLog_Append)
}

func testEventLog_Append(t *testing.T) {
	ctx := t.Context()
	el := New(createStorage(t, 2, 1024))
	fn := path.Join(t.TempDir(), "testfile")

	records := []*entity.Event{{
		GlobalVersion: value.NewGlobalVersion(0),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(1),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(2),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(3),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(4),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}}

	for _, r := range records {
		if _, err := el.Append(ctx, fn, r); err != nil {
			t.Fatal(err)
		}
	}

	result := make([]*entity.Event, 0)

	for e, err := range el.Events(ctx, fn) {
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
	el := New(createStorage(t, 2, 1024))
	fn := path.Join(t.TempDir(), "testfile")

	records := []*entity.Event{{
		GlobalVersion: value.NewGlobalVersion(0),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(1),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(2),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(3),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(4),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}}

	for _, r := range records {
		if _, err := el.Append(ctx, fn, r); err != nil {
			t.Fatal(err)
		}
	}

	e, err := el.Find(ctx, fn, records[2].ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(e, records[2]) {
		t.Errorf("incorrect event: %v, expected: %v", e, records[2])
	}
}

func TestEventLog_GetAggregate(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventLog_GetAggregate)
}

func testEventLog_GetAggregate(t *testing.T) {
	ctx := t.Context()
	el := New(createStorage(t, 2, 1024))
	fn := path.Join(t.TempDir(), "testfile")

	aggregateID := value.NewAggregateID(uuid.NewString())

	records := []*entity.Event{{
		GlobalVersion: value.NewGlobalVersion(0),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(1),
		ID:            value.NewEventID(),
		AggregateID:   aggregateID,
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(2),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(3),
		ID:            value.NewEventID(),
		AggregateID:   aggregateID,
		Data:          value.NewData([]byte{}),
	}, {
		GlobalVersion: value.NewGlobalVersion(4),
		ID:            value.NewEventID(),
		Data:          value.NewData([]byte{}),
	}}

	for _, r := range records {
		if _, err := el.Append(ctx, fn, r); err != nil {
			t.Fatal(err)
		}
	}

	want := []*entity.Event{records[1], records[3]}

	result := make([]*entity.Event, 0)

	for e, err := range el.GetAggregate(ctx, fn, aggregateID) {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, e)
	}

	if !slices.EqualFunc(result, want, func(a, b *entity.Event) bool {
		return reflect.DeepEqual(a, b)
	}) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func TestEventLog_AppendEvent(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventLog_AppendEvent)
}

func testEventLog_AppendEvent(t *testing.T) {
	ctx := t.Context()
	el := New(createStorage(t, 1, 48))
	fn := path.Join(t.TempDir(), "testfile")

	records := [][]byte{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10}}

	for _, r := range records {
		if _, err := el.AppendEvent(ctx, fn, raw.NewBufferFromSlice(r)); err != nil {
			t.Fatal(err)
		}
	}

	pageCount := 0
	result := make([][]byte, 0)

	for b, err := range el.Iterate(ctx, fn) {
		if err != nil {
			t.Fatal(err)
		}

		pageCount++

		for i, err := range b.Items() {
			if err != nil {
				t.Fatal(err)
			}

			result = append(result, i)
		}
	}

	if pageCount != 3 {
		t.Errorf("incorrect page count: %v, expected: %v", pageCount, 3)
	}

	if !slices.EqualFunc(result, records, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, records)
	}
}
