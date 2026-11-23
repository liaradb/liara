package eventlog

import (
	"context"
	"reflect"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/heap"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/infrastructure/tenantfolder"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/lrupool"
	"github.com/liaradb/liaradb/util/testing/filetesting"
	"github.com/liaradb/liaradb/util/testing/txtesting"
)

func TestEventLog_Recovery(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testRecovery)
}

func testRecovery(t *testing.T) {
	ctx := t.Context()

	fsys := filetesting.New(nil)
	dir := t.TempDir()
	tf := tenantfolder.NewTenantFolder(value.NewTenantID())
	pid := value.NewPartitionID(0)
	tn := tf.EventLog(pid)

	var max int = 2
	var bs int16 = 1024

	records := []*entity.Event{
		createEvent(0, 0),
		createEvent(1, 1),
		createEvent(2, 2),
		createEvent(3, 3),
		createEvent(4, 4)}

	write(t, ctx, fsys, max, bs, dir, tn, records)
	recover(t, ctx, fsys, max, bs, dir, tn, records)
}

func write(
	t *testing.T,
	baseCtx context.Context,
	fsys filecache.FileSystem,
	max int,
	bs int16,
	dir string,
	tn tablename.TableName,
	events []*entity.Event,
) {
	s := storage.New(fsys, lrupool.New(), max, bs, dir)
	el := New(indexed.New(heap.New()))

	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()

	if err := s.Run(ctx); err != nil {
		t.Fatal(err)
	}

	tx := txtesting.New(s)

	for _, r := range events {
		if err := el.Append(ctx, tx, tn, r); err != nil {
			t.Fatal(err)
		}
	}

	tx.Release()

	if err := s.FlushUnpinned(baseCtx); err != nil {
		t.Fatal(err)
	}
}

func recover(
	t *testing.T,
	ctx context.Context,
	fsys filecache.FileSystem,
	max int,
	bs int16,
	dir string,
	tn tablename.TableName,
	events []*entity.Event,
) {
	s := storage.New(fsys, lrupool.New(), max, bs, dir)
	el := New(indexed.New(heap.New()))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := s.Run(ctx); err != nil {
		t.Fatal(err)
	}

	result := make([]*entity.Event, 0)

	tx := txtesting.New(s)

	defer tx.Release()

	for e, err := range el.EventsAfterGlobalVersion(ctx, tx, tn,
		value.NewGlobalVersion(0),
		value.NewSubVersion(0),
		value.NewGlobalVersion(10)) {
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, e)
	}

	if !slices.EqualFunc(result, events, func(a, b *entity.Event) bool {
		return reflect.DeepEqual(a, b)
	}) {
		t.Errorf("incorrect result:\n%v\nexpected:\n%v", result, events)
	}
}
