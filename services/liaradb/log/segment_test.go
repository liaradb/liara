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
		names = append(names, NewSegmentName(i, 0))
	}
	return names
}

func createFiles(count SegmentID) fs.FS {
	fsys := fstest.MapFS{}
	for i := range count {
		fsys[NewSegmentName(i, 0).String()] = &fstest.MapFile{}
	}
	return fsys
}
