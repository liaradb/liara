package log

import (
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestListSegments(t *testing.T) {
	t.Parallel()

	count := 10

	fsys := createFiles(count)
	names, err := ListSegments(".", fsys)
	if err != nil {
		t.Fatal(err)
	}

	want := createNames(count)
	if !reflect.DeepEqual(want, names) {
		t.Errorf("files do not match:\n\t%v,\nexpected:\n\t%v", names, want)
	}
}

func createNames(count int) []LogSegmentName {
	names := make([]LogSegmentName, 0, count)
	for i := range count {
		names = append(names, NewLogSegmentName(i, 0))
	}
	return names
}

func createFiles(count int) fs.FS {
	fsys := fstest.MapFS{}
	for i := range count {
		fsys[NewLogSegmentName(i, 0).String()] = &fstest.MapFile{}
	}
	return fsys
}
