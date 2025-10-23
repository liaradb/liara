package eventlog

import (
	"bytes"
	"path"
	"reflect"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/raw"
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
		if err := el.Append(ctx, fn, r); err != nil {
			t.Fatal(err)
		}
	}

	result := make([]*entity.Event, 0)

	buf := bytes.NewBuffer(nil)

	for b, err := range el.Iterate(ctx, fn) {
		if err != nil {
			t.Fatal(err)
		}

		// pageCount++

		for i, err := range b.Items() {
			if err != nil {
				t.Fatal(err)
			}

			buf = bytes.NewBuffer(i)

			var e entity.Event
			e.Read(buf)

			result = append(result, &e)
		}
	}

	if !slices.EqualFunc(result, records, func(a, b *entity.Event) bool {
		return reflect.DeepEqual(a, b)
	}) {
		t.Errorf("incorrect result: %v, expected: %v", result, records)
	}
}

func TestEventLog_AppendEvent(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventLog_AppendEvent)
}

func testEventLog_AppendEvent(t *testing.T) {
	ctx := t.Context()
	el := New(createStorage(t, 1, 32))
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
