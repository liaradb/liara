package log

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func TestListSegments(t *testing.T) {
	t.Parallel()

	count := 10

	dir := t.TempDir()
	if err := createFiles(dir, count); err != nil {
		t.Fatal(err)
	}

	names, err := ListSegments(dir)
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
		names = append(names, NewLogSegmentName(i))
	}
	return names
}

func createFiles(dir string, count int) error {
	for i := range count {
		if err := createFile(dir, i); err != nil {
			return err
		}
	}
	return nil
}

func createFile(dir string, index int) error {
	name := path.Join(dir, NewLogSegmentName(index).String())
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	return f.Close()
}
