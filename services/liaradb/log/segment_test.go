package log

import (
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestSegment(t *testing.T) {
	t.Parallel()

	t.Run("should handle default", func(t *testing.T) {
		t.Parallel()
		s := NewSegment(0, 0)

		if v := s.Size(); v != 0 {
			t.Errorf("incorrect size: %v, expected: %v", v, 0)
		}

		if v := s.PageSize(); v != 0 {
			t.Errorf("incorrect page size: %v, expected: %v", v, 0)
		}
	})

	t.Run("should handle values", func(t *testing.T) {
		t.Parallel()
		s := NewSegment(1, 2)

		if v := s.Size(); v != 1 {
			t.Errorf("incorrect size: %v, expected: %v", v, 1)
		}

		if v := s.PageSize(); v != 2 {
			t.Errorf("incorrect page size: %v, expected: %v", v, 2)
		}
	})
}

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
		names = append(names, NewSegmentName(i, 0))
	}
	return names
}

func createFiles(count SegmentID) fs.ReadDirFS {
	fsys := fstest.MapFS{}
	for i := range count {
		fsys[NewSegmentName(i, 0).String()] = &fstest.MapFile{}
	}
	return fsys
}
