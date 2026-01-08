package eventlog

import (
	"context"
	"path"
	"reflect"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/file/filetesting"
	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/storage"
)

// TODO: Where should this test be?
// TODO: Why does this take so long?
func TestEventLog_Recovery(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testRecovery)
}

func testRecovery(t *testing.T) {
	ctx := t.Context()

	fsys := filetesting.NewMockFileSystem(t, nil)
	dir := t.TempDir()
	tn := tablename.New(value.TenantID(path.Join(dir, "testfile")))
	pid := value.NewPartitionID(0)

	var max int = 2
	var bs int64 = 256

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

	write(t, ctx, fsys, max, bs, dir, tn, pid, records)
	recover(t, ctx, fsys, max, bs, dir, tn, pid, records)
}

func write(
	t *testing.T,
	baseCtx context.Context,
	fsys *mock.FileSystem,
	max int,
	bs int64,
	dir string,
	tn tablename.TableName,
	pid value.PartitionID,
	records []*entity.Event,
) {
	s := storage.New(fsys, max, bs, dir)
	el := New(s)

	ctx, cancel := context.WithCancel(baseCtx)

	if err := s.Run(ctx); err != nil {
		t.Fatal(err)
	}

	for _, r := range records {
		if _, err := el.Append(ctx, tn, pid, r); err != nil {
			t.Fatal(err)
		}
	}

	cancel()

	if err := s.FlushAll(); err != nil {
		t.Fatal(err)
	}
}

func recover(
	t *testing.T,
	ctx context.Context,
	fsys *mock.FileSystem,
	max int,
	bs int64,
	dir string,
	tn tablename.TableName,
	pid value.PartitionID,
	records []*entity.Event,
) {
	s := storage.New(fsys, max, bs, dir)
	el := New(s)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := s.Run(ctx); err != nil {
		t.Fatal(err)
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
		t.Errorf("incorrect result:\n%v\nexpected:\n%v", result, records)
	}
}
