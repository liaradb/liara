package eventlog

import (
	"path"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/raw"
	"github.com/liaradb/liaradb/storage"
)

func TestEventLog_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventLog_AppendEvent)
}

func testEventLog_AppendEvent(t *testing.T) {
	ctx := t.Context()
	s := createStorage(t, 1, 32)
	el := New(s)
	n := path.Join(t.TempDir(), "testfile")

	records := [][]byte{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10},
	}

	(func() {
		var bid storage.BlockID
		var err error
		for _, r := range records {
			if bid, err = el.AppendEvent(ctx, n, raw.NewBufferFromSlice(r)); err != nil {
				t.Error(err)
			}
			synctest.Wait()
		}

		b, err := s.Request(ctx, bid)
		if err != nil {
			t.Fatal(err)
		}

		defer b.Release()

		err = b.Flush()
		if err != nil {
			t.Fatal(err)
		}
	})()

	pageCount := 0
	recordCount := 0
	result := make([][]byte, 0, len(records))
	for b, err := range el.Iterate(ctx, n) {
		if err != nil {
			t.Fatal(err)
		}

		pageCount++

		for i, err := range b.Items() {
			if err != nil {
				t.Fatal(err)
			}

			recordCount++

			result = append(result, i)
		}
	}

	if pageCount != 3 {
		t.Errorf("incorrect page count: %v, expected: %v", pageCount, 3)
	}

	if recordCount != len(records) {
		t.Fatalf("incorrect record count: %v, expected: %v", recordCount, len(records))
	}

	for index, r := range records {
		i := result[index]
		if !slices.Equal(i, r) {
			t.Errorf("incorrect record: %v, expected: %v", i, r)
		}
	}
}
