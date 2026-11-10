package tip

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/logtesting"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

const (
	pageSize       = 78
	largePageSize  = 256
	writeQueueSize = 100
)

func TestTip(t *testing.T) {
	storagetesting.SyncTest(t, 16, pageSize, func(t *testing.T, st storagetesting.Storage) {
		tip := NewTip(st.Storage, link.NewFileName("fn"))
		want := 128
		s, err := tip.Span(t.Context(), want)
		if err != nil {
			t.Fatal(err)
		}

		if l := s.Length(); l != want {
			t.Errorf("incorrect length: %v, expected: %v", l, want)
		}

		l := &logtesting.MockLog{}
		if err := s.CommitAppend(t.Context(), l); err != nil {
			t.Fatal(err)
		}

		s.Release()

		size := 0
		for _, r := range l.Records() {
			size += len(r.Data())
		}
		if size != want {
			t.Errorf("incorrect length: %v, expected: %v", l, want)
		}
		if c := len(l.Records()); c != 3 {
			t.Errorf("incorrect complete count: %v, expected: %v", c, 3)
		}

		synctest.Wait()
	})
}
