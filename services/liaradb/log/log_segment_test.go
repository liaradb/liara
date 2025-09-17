package log

import (
	"errors"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestListSegments(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	for i := range 10 {
		if err := createFile(dir, i); err != nil {
			t.Fatal(err)
		}
	}

	names, err := ListSegments(dir)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(createNames(10), names) {
		t.Error("files to not match")
	}
}

func createNames(count int) []LogSegmentName {
	names := make([]LogSegmentName, 0, count)
	for i := range count {
		names = append(names, NewLogSegmentName(i))
	}
	return names
}

func createFile(dir string, index int) (err error) {
	name := path.Join(dir, NewLogSegmentName(index).String())
	var f *os.File
	f, err = os.Create(name)
	if err != nil {
		return
	}

	defer func() { errors.Join(err, f.Close()) }()

	return
}
