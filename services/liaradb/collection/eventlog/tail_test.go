package eventlog

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/storage/storagetesting"
)

func TestTail(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTail)
}

func testTail(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 16)
	tl := NewTail(s)
	if err := tl.Append(t.Context()); err != nil {
		t.Error(err)
	}
}
