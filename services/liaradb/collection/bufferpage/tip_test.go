package bufferpage

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

const (
	pageSize       = 64
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

		// complete := 0
		pages, ok := tip.Commit()
		if !ok {
			t.Error("should commit")
		}

		if l := len(pages); l != 3 {
			t.Errorf("incorrect length: %v, expected: %v", l, 3)
		}

		for _, p := range pages {
			p.Release()
		}

		// last := pages[len(pages)-1]
		// last.Complete()
		// if complete != 1 {
		// 	t.Errorf("incorrect complete count: %v, expected: %v", complete, 1)
		// }
		synctest.Wait()
	})
}
