package manager

import (
	"slices"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/storage/storagetesting"
)

func TestManager(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testManager)
}

func testManager(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	m := New(s)

	type tuple struct {
		key   value.Key
		value int64
	}
	data := []tuple{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 4},
		{"e", 5}}

	d := slices.Clone(data)
	slices.SortFunc(d, func(a, b tuple) int {
		return strings.Compare(a.key.String(), b.key.String())
	})
	want := make([]int64, 0, len(d))
	for _, d := range data {
		want = append(want, d.value)
	}

	for _, d := range data {
		if err := m.Insert(t.Context(), d.key, d.value); err != nil {
			t.Fatal(err)
		}
	}

	i := make([]int64, 0, len(d))
	for _, d := range data {
		v, err := m.Get(t.Context(), d.key)
		if err != nil {
			t.Fatal(err)
		}

		i = append(i, v)
	}

	if !slices.Equal(i, want) {
		t.Errorf("incorrect result: %v, expected: %v", i, want)
	}
}
