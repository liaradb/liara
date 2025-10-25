package eventlog

import (
	"context"
	"path"
	"reflect"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filetesting"
	"github.com/liaradb/liaradb/storage"
)

func TestEventLog_Recovery(t *testing.T) {
	t.Parallel()
	t.Skip()
	synctest.Test(t, testRecovery)
}

func testRecovery(t *testing.T) {
	baseCtx := t.Context()

	fsys := filetesting.NewMockFileSystem(t, nil)
	dir := t.TempDir()
	fn := path.Join(dir, "testfile")
	var max int = 2
	var bs int64 = 1024

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

	{
		s := storage.New(fsys, max, bs, dir)
		el := New(s)

		ctx, cancel := context.WithCancel(baseCtx)

		if err := s.Run(ctx); err != nil {
			t.Fatal(err)
		}

		for _, r := range records {
			if err := el.Append(ctx, fn, r); err != nil {
				t.Fatal(err)
			}
		}

		cancel()

		if err := s.FlushAll(); err != nil {
			t.Fatal(err)
		}
	}
	{
		s := storage.New(fsys, max, bs, dir)
		el := New(s)

		ctx, cancel := context.WithCancel(baseCtx)
		defer cancel()

		if err := s.Run(ctx); err != nil {
			t.Fatal(err)
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
}
