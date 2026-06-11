package pagestorage

import (
	"io/fs"
	"path"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/util/testing/filetesting"
)

func TestPageStorage(t *testing.T) {
	t.Parallel()

	fsys := filetesting.New(nil)
	dir := "dir"
	sl := segment.NewList(fsys, dir, 128, 1)
	ps := New(sl)

	pageData := make([]byte, 128)
	if err := ps.Init(pageData); err != nil {
		t.Error(err)
	}

	page0 := make([]byte, 128)
	for i := range len(page0) {
		page0[i] = 1
	}

	if err := ps.Sync(page0); err != nil {
		t.Error(err)
	}

	page1 := make([]byte, 128)
	for i := range len(page1) {
		page1[i] = 2
	}

	if err := ps.Append(record.NewLogSequenceNumber(1), page1); err != nil {
		t.Error(err)
	}

	if err := sl.Close(); err != nil {
		t.Error(err)
	}

	files, err := fsys.ReadDir(dir)
	if err != nil {
		t.Error(err)
	}

	if c := len(files); c != 2 {
		t.Errorf("incorrect count: %v, expected: %v", c, 2)
	}

	for i, de := range files {
		switch i {
		case 0:
			testPage(t, de, fsys, dir, page0)
		case 1:
			testPage(t, de, fsys, dir, page1)
		}
	}
}

func testPage(t *testing.T, de fs.DirEntry, fsys *filecache.Cache, dir string, page []byte) {
	fi, err := de.Info()
	if err != nil {
		t.Fatal(err)
	}
	if fi.Size() != 128 {
		t.Errorf("incorrect size: %v, expected: %v", fi.Size(), 128)
	}
	f, err := fsys.OpenFile(path.Join(dir, fi.Name()))
	if err != nil {
		t.Fatal(err)
	}

	result := make([]byte, 128)
	if _, err := f.Read(result); err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(result, page) {
		t.Errorf("incorrect data: %v, expected: %v", result, page)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}
