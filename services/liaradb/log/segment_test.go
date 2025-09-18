package log

import (
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestGetSegmentForLSN(t *testing.T) {
	fsys := fstest.MapFS{
		NewLogSegmentName(1, 10).String(): {},
		NewLogSegmentName(2, 20).String(): {},
	}
	names, err := ListSegments(fsys, ".")
	if err != nil {
		t.Fatal(err)
	}

	for message, test := range map[string]struct {
		search LogSequenceNumber
		found  bool
		result LogSequenceNumber
	}{
		"should not find low value": {1, false, 0},
		"should find exact value":   {10, true, 10},
		"should find middle value":  {15, true, 10},
		"should find high value":    {50, true, 20},
	} {
		t.Run(message, func(t *testing.T) {
			lsn, ok := GetSegmentForLSN(names, test.search)
			if test.found {
				if !ok {
					t.Error("should find LSN")
				}
				if lsn.lsn != test.result {
					t.Error("wrong LSN")
				}
			} else {
				if ok {
					t.Error("should not find LSN")
				}
			}
		})
	}
}

func TestListSegments(t *testing.T) {
	t.Parallel()

	var count SegmentID = 10

	fsys := createFiles(count)
	names, err := ListSegments(fsys, ".")
	if err != nil {
		t.Fatal(err)
	}

	want := createNames(count)
	if !reflect.DeepEqual(want, names) {
		t.Errorf("files do not match:\n\t%v,\nexpected:\n\t%v", names, want)
	}
}

func createNames(count SegmentID) []SegmentName {
	names := make([]SegmentName, 0, count)
	for i := range count {
		names = append(names, NewLogSegmentName(i, 0))
	}
	return names
}

func createFiles(count SegmentID) fs.FS {
	fsys := fstest.MapFS{}
	for i := range count {
		fsys[NewLogSegmentName(i, 0).String()] = &fstest.MapFile{}
	}
	return fsys
}
