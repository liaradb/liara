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
	t.Skip()
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
		// fmt.Println(s.pinned)
		var bid storage.BlockID
		var err error
		for _, r := range records {
			if _, err = el.AppendEvent(ctx, n, raw.NewBufferFromSlice(r)); err != nil {
				t.Error(err)
			}
			synctest.Wait()
			// fmt.Println(s.pinned)
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

	c := 0
	l := len(records) - 1
	for b, err := range el.Iterate(ctx, n) {
		if err != nil {
			t.Error(err)
		}

		for i, err := range b.Items() {
			if err != nil {
				t.Error(err)
			}

			// fmt.Println(i)

			r := records[l-c]
			if !slices.Equal(i, r) {
				t.Errorf("incorrect record: %v, expected: %v", i, r)
			}

			c++
		}
	}

	if c != len(records) {
		t.Errorf("incorrect count: %v, expected: %v", c, len(records))
	}
}
