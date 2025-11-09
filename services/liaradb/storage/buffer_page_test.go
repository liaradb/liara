package storage

import (
	"path"
	"slices"
	"testing"
	"testing/synctest"
)

func TestBufferPage(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testBufferPage)
}

func testBufferPage(t *testing.T) {
	b := testCreateBufferPage(t)

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	want := [][]byte{{1, 2, 3, 4, 5}}

	if err := b.Add(want[0]); err != nil {
		t.Fatal(err)
	}

	if !b.Dirty() {
		t.Error("should be dirty")
	}

	if err := b.Flush(); err != nil {
		t.Fatal(err)
	}

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	if err := b.Flush(); err == nil {
		t.Fatal("should not flush clean buffers")
	}

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	result := make([][]byte, 0)

	for i, err := range b.Items() {
		if err != nil {
			t.Error(err)
		}

		result = append(result, i)
	}

	if !slices.EqualFunc(result, want, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func testCreateBufferPage(t *testing.T) *BufferPage {
	b, err := createStorage(t, 2, 1024).
		Request(t.Context(),
			NewBlockID(path.Join(t.TempDir(), "testfile"), 0))
	if err != nil {
		t.Fatal(err)
	}

	return NewBufferPage(b)
}
