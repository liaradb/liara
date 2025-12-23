package manager

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/storage/storagetesting"
)

func TestManager(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testManager)
}

func testManager(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	m := New(s)

	var want int64 = 2
	if err := m.Insert(t.Context(), "a", want); err != nil {
		t.Fatal(err)
	}

	i, err := m.List(t.Context(), "a")
	if err != nil {
		t.Fatal(err)
	}

	if i != want {
		t.Errorf("incorrect result: %v, expected: %v", i, want)
	}
}
