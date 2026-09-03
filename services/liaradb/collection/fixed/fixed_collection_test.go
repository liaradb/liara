package fixed

import (
	"context"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/google/uuid"
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestFixedCollection_InsertAndGet(t *testing.T) {
	storagetesting.SyncTest(t, 16, 1024, func(t *testing.T, s storagetesting.Storage) {
		l := log.New(256, 2, 256, 100, s.FSys, "dir")
		if err := l.Run(t.Context()); err != nil {
			t.Fatal(err)
		}

		fc := New(s.Storage, btree.NewCursor(s.Storage), l)

		fn := link.NewFileName("testfile")
		fnIdx := link.NewFileName("testindex")
		// pid := value.NewPartitionID(0)
		k := key.NewKey([]byte("abcde"))
		want := []byte{1, 2, 3, 4, 5}

		if err := fc.Insert(t.Context(), fn, fnIdx, k, want); err != nil {
			t.Fatal(err)
		}

		result, err := fc.Get(t.Context(), fn, fnIdx, k)
		if err != nil {
			t.Fatal(err)
		}

		if !slices.Equal(result, want) {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	})
}

func TestFixedCollection(t *testing.T) {
	storagetesting.SyncTest(t, 20, 110, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		l := log.New(256, 2, 256, 100, s.FSys, "dir")
		if err := l.Run(t.Context()); err != nil {
			t.Fatal(err)
		}

		fc := New(s.Storage, btree.NewCursor(s.Storage), l)
		fn := link.NewFileName("testfile")
		fnIdx := link.NewFileName("testindex")
		pid := value.NewPartitionID(0)

		data := createData()
		slices.Reverse(data)

		if err := insertData(ctx, fc, fn, fnIdx, data); err != nil {
			t.Fatal(err)
		}

		testGet(ctx, t, fc, fn, fnIdx, data)
		testList(ctx, t, data, fc, fn, fnIdx, pid)

		synctest.Wait()
	})
}

func TestRequestLog__LargeBuffer(t *testing.T) {
	t.Skip()
	storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		l := log.New(256, 2, 256, 100, s.FSys, "dir")
		if err := l.Run(t.Context()); err != nil {
			t.Fatal(err)
		}

		fc := New(s.Storage, btree.NewCursor(s.Storage), l)
		fn := link.NewFileName("testfile")
		fnIdx := link.NewFileName("testindex")
		pid := value.NewPartitionID(0)

		data := createData()

		if err := insertData(ctx, fc, fn, fnIdx, data); err != nil {
			t.Fatal(err)
		}

		testGet(ctx, t, fc, fn, fnIdx, data)
		testList(ctx, t, data, fc, fn, fnIdx, pid)

		synctest.Wait()
	})
}

type item struct {
	key   string
	value *entity.RequestLog
}

func createData() []item {
	count := 9
	data := uuid.UUID{}
	items := make([]item, 0, count)
	for i := range count {
		data[15] = byte(i) + 1
		rid, _ := value.NewRequestIDFromString(data.String())
		items = append(items, item{fmt.Sprintf("%v", i+1), entity.NewRequestLog(rid, value.NewTime(time.Now().Add(time.Duration(i)*time.Second)))})
	}
	return items
}

func insertData(ctx context.Context, fc *FixedCollection, fn link.FileName, fnIdx link.FileName, data []item) error {
	for _, i := range data {
		d := make([]byte, entity.RequestLogSize)

		if _, ok := i.value.Write(d); !ok {
			return io.EOF
		}

		k := key.NewKey(i.value.ID().Bytes())
		if err := fc.Insert(ctx, fn, fnIdx, k, d); err != nil {
			return err
		}
	}
	return nil
}

func testGet(
	ctx context.Context,
	t *testing.T,
	fc *FixedCollection,
	fn link.FileName,
	fnIdx link.FileName,
	data []item,
) {
	for _, i := range data {
		k := key.NewKey(i.value.ID().Bytes())
		value, err := fc.Get(ctx, fn, fnIdx, k)
		if err != nil {
			t.Fatal(i.key, err)
		}

		rl := entity.RequestLog{}
		if _, ok := rl.Read(value); !ok {
			t.Fatal("should read")
		}

		if rl != *i.value {
			t.Errorf("incorrect result: %v, expected: %v", rl, *i.value)
		}
	}
}

func testList(
	ctx context.Context,
	t *testing.T,
	data []item,
	fc *FixedCollection,
	fn link.FileName,
	fnIdx link.FileName,
	pid value.PartitionID,
) {
	result, err := getListValues(ctx, data, fc, fn, fnIdx, pid)
	if err != nil {
		t.Fatal(err)
	}

	want := createSortedValues(data)
	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func getListValues(
	ctx context.Context,
	data []item,
	fc *FixedCollection,
	fn link.FileName,
	fnIdx link.FileName,
	pid value.PartitionID,
) ([]entity.RequestLog, error) {
	result := make([]entity.RequestLog, 0, len(data))
	i := 0
	for value, err := range fc.List(ctx, fn, fnIdx, pid) {
		if err != nil {
			return nil, err
		}

		rl := entity.RequestLog{}
		if _, ok := rl.Read(value); !ok {
			return nil, io.EOF
		}

		result = append(result, rl)
		i++
	}
	return result, nil
}

func createSortedValues(data []item) []entity.RequestLog {
	type tuple struct {
		key   key.Key
		value *entity.RequestLog
	}

	tuples := make([]tuple, 0, len(data))
	for _, i := range data {
		tuples = append(tuples, tuple{key.NewKey(i.value.ID().Bytes()), i.value})
	}
	slices.SortFunc(tuples, func(a, b tuple) int {
		return strings.Compare(a.key.String(), b.key.String())
	})
	want := make([]entity.RequestLog, 0, len(data))
	for _, t := range tuples {
		want = append(want, *t.value)
	}
	return want
}
