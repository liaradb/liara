package log

import (
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestGetLatestSegment(t *testing.T) {
	t.Parallel()

	for message, test := range map[string]struct {
		result SegmentID
		fsys   fs.ReadDirFS
	}{
		"should handle no files": {0, fstest.MapFS{}},
		"should handle one file": {1, fstest.MapFS{
			NewSegmentName(1, 10).String(): {},
		}},
		"should handle multiple files": {2, fstest.MapFS{
			NewSegmentName(1, 10).String(): {},
			NewSegmentName(2, 20).String(): {},
		}},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()

			names, err := ListSegments(test.fsys, ".")
			if err != nil {
				t.Fatal(err)
			}

			sn := GetLatestSegment(names)
			if id := sn.ID(); id != test.result {
				t.Errorf("wrong id: %v, expected: %v", id, test.result)
			}
		})
	}
}

func TestGetSegmentForLSN(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		NewSegmentName(1, 10).String(): {},
		NewSegmentName(2, 20).String(): {},
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
			sn, ok := GetSegmentForLSN(names, test.search)
			if test.found {
				if !ok {
					t.Error("should find log sequence number")
				}
				if lsn := sn.LogSequenceNumber(); lsn != test.result {
					t.Errorf("wrong log sequence number: %v, expected: %v", lsn, test.result)
				}
			} else {
				if ok {
					t.Error("should not find log sequence number")
				}
			}
		})
	}
}

func TestListSegments(t *testing.T) {
	t.Parallel()

	var count SegmentID = 10

	t.Run("should list segments", func(t *testing.T) {
		fsys := createFiles(0, count)
		names, err := ListSegments(fsys, ".")
		if err != nil {
			t.Fatal(err)
		}

		want := createNames(0, count)
		if !reflect.DeepEqual(want, names) {
			t.Errorf("files do not match:\n\t%v,\nexpected:\n\t%v", names, want)
		}
	})

	t.Run("should list segments in order", func(t *testing.T) {
		fsys := createFiles(9998, count)
		names, err := ListSegments(fsys, ".")
		if err != nil {
			t.Fatal(err)
		}

		want := createNames(9998, count)
		if !reflect.DeepEqual(want, names) {
			t.Errorf("files do not match:\n\t%v,\nexpected:\n\t%v", names, want)
		}
	})
}

func createNames(start SegmentID, count SegmentID) []SegmentName {
	names := make([]SegmentName, 0, count)
	for i := range count {
		names = append(names, NewSegmentName(start+i, 0))
	}
	return names
}

func createFiles(start SegmentID, count SegmentID) fs.ReadDirFS {
	fsys := fstest.MapFS{}
	for i := range count {
		fsys[NewSegmentName(start+i, 0).String()] = &fstest.MapFile{}
	}
	return fsys
}
