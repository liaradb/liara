package manager

import (
	"slices"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/storage/storagetesting"
)

func TestManager(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testManager)
}

func testManager(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	m := New(keyvalue.New(s, btree.NewCursor(s)))

	data := createData()
	want := createValues(data)

	for _, d := range data {
		if err := m.Insert(t.Context(), value.NewKey([]byte(d.key)), d.value); err != nil {
			t.Fatal(err)
		}
	}

	testGet(t, data, want, m)
	testList(t, want, m)

	synctest.Wait()
}

func testGet(t *testing.T, data []tuple, want []int64, m *Manager) {
	i := make([]int64, 0, len(data))
	for _, d := range data {
		v, err := m.Get(t.Context(), value.NewKey([]byte(d.key)))
		if err != nil {
			t.Fatal(err)
		}

		i = append(i, v)
	}

	if !slices.Equal(i, want) {
		t.Errorf("incorrect result: %v, expected: %v", i, want)
	}
}

func testList(t *testing.T, want []int64, m *Manager) {
	i, err := m.List(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(i, want) {
		t.Errorf("incorrect result: %v, expected: %v", i, want)
	}
}

type tuple struct {
	key   string
	value int64
}

func createData() []tuple {
	return []tuple{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 4},
		{"e", 5}}
}

func createValues(data []tuple) []int64 {
	d := slices.Clone(data)
	slices.SortFunc(d, func(a, b tuple) int {
		return strings.Compare(a.key, b.key)
	})
	values := make([]int64, 0, len(d))
	for _, d := range data {
		values = append(values, d.value)
	}
	return values
}
