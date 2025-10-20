package eventlog

import (
	"context"
	"path"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/raw"
	"github.com/liaradb/liaradb/storage"
)

func TestEventLog_AppendEvent(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testEventLog_AppendEvent)
}

func testEventLog_AppendEvent(t *testing.T) {
	ctx := t.Context()
	s := createStorage(t, 1, 32)
	el := New(s)
	fn := path.Join(t.TempDir(), "testfile")

	records := [][]byte{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10},
	}

	if err := appendRecords(ctx, s, el, fn, records); err != nil {
		t.Fatal(err)
	}

	pageCount, recordCount, result, err := iterateRecords(ctx, el, fn)
	if err != nil {
		t.Fatal(err)
	}

	testRecords(t, pageCount, recordCount, records, result)
}

func appendRecords(ctx context.Context, s *storage.Storage, el *EventLog, fn string, records [][]byte) error {
	var bid storage.BlockID
	var err error
	for _, r := range records {
		if bid, err = el.AppendEvent(ctx, fn, raw.NewBufferFromSlice(r)); err != nil {
			return err
		}
		synctest.Wait()
	}

	b, err := s.Request(ctx, bid)
	if err != nil {
		return err
	}

	defer b.Release()

	return b.Flush()
}

func iterateRecords(ctx context.Context, el *EventLog, fn string) (int, int, [][]byte, error) {
	pageCount := 0
	recordCount := 0
	result := make([][]byte, 0)

	for b, err := range el.Iterate(ctx, fn) {
		if err != nil {
			return 0, 0, nil, err
		}

		pageCount++

		for i, err := range b.Items() {
			if err != nil {
				return 0, 0, nil, err
			}

			recordCount++

			result = append(result, i)
		}
	}

	return pageCount, recordCount, result, nil
}

func testRecords(t *testing.T, pageCount int, recordCount int, records [][]byte, result [][]byte) {
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
