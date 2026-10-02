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
		tip := NewTip(st.Storage, &logtesting.MockLog{}, link.NewFileName("fn"))
		want := 128
		s, err := tip.Span(t.Context(), want)
		if err != nil {
			t.Fatal(err)
		}

		if l := s.Length(); l != want {
			t.Errorf("incorrect length: %v, expected: %v", l, want)
		}

		// complete := 0
		// TODO: Test the number of pages
		tip.Commit()

		tip.Release()

		// last := pages[len(pages)-1]
		// last.Complete()
		// if complete != 1 {
		// 	t.Errorf("incorrect complete count: %v, expected: %v", complete, 1)
		// }
		synctest.Wait()
	})
}
