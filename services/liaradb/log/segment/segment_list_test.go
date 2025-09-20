package segment

import (
	"reflect"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/log/record"
)

func TestSegmentList_Open(t *testing.T) {
	t.Parallel()

	var count SegmentID = 10

	t.Run("should list segments", func(t *testing.T) {
		t.Parallel()

		fsys := createFiles(0, count)
		sl := NewSegmentList(fsys, ".")

		if err := sl.Open(); err != nil {
			t.Fatal(err)
		}

		names := sl.Names()

		want := createNames(0, count)
		if !reflect.DeepEqual(want, names) {
			t.Errorf("files do not match:\n\t%v,\nexpected:\n\t%v", names, want)
		}
	})

	t.Run("should list segments in order", func(t *testing.T) {
		t.Parallel()

		fsys := createFiles(9998, count)
		sl := NewSegmentList(fsys, ".")

		if err := sl.Open(); err != nil {
			t.Fatal(err)
		}

		names := sl.Names()

		want := createNames(9998, count)
		if !reflect.DeepEqual(want, names) {
			t.Errorf("files do not match:\n\t%v,\nexpected:\n\t%v", names, want)
		}
	})
}

func TestSegmentList_OpenLatestSegment(t *testing.T) {
	t.Parallel()

	for message, test := range map[string]struct {
		result SegmentID
		fsys   file.FileSystem
	}{
		"should handle no files": {0, mock.NewFileSystem(fstest.MapFS{})},
		"should handle one file": {1, mock.NewFileSystem(fstest.MapFS{
			NewSegmentName(1, 10).String(): {},
		})},
		"should handle multiple files": {2, mock.NewFileSystem(fstest.MapFS{
			NewSegmentName(1, 10).String(): {},
			NewSegmentName(2, 20).String(): {},
		})},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()

			sl := NewSegmentList(test.fsys, ".")

			if err := sl.Open(); err != nil {
				t.Fatal(err)
			}

			sn, f, err := sl.OpenLatestSegment()
			if err != nil {
				t.Fatal(err)
			}

			if id := sn.ID(); id != test.result {
				t.Errorf("wrong id: %v, expected: %v", id, test.result)
			}

			if f == nil {
				t.Error("file should not be nil")
			}

			if names := sl.Names(); !slices.Contains(names, sn) {
				t.Errorf("segment list does not contain segment: %v", sn)
			}
		})
	}
}

func TestSegmentList_OpenSegmentForLSN(t *testing.T) {
	t.Parallel()

	fsys := mock.NewFileSystem(fstest.MapFS{
		NewSegmentName(1, 10).String(): {},
		NewSegmentName(2, 20).String(): {},
	})
	sl := NewSegmentList(fsys, ".")

	if err := sl.Open(); err != nil {
		t.Fatal(err)
	}

	for message, test := range map[string]struct {
		search record.LogSequenceNumber
		found  bool
		result record.LogSequenceNumber
	}{
		"should not find low value": {1, false, 0},
		"should find exact value":   {10, true, 10},
		"should find middle value":  {15, true, 10},
		"should find high value":    {50, true, 20},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()

			sn, f, err := sl.OpenSegmentForLSN(test.search)
			if test.found {
				if err != nil {
					if err == ErrNoSegmentFile {
						t.Error("should find log sequence number")
					} else {
						t.Error(err)
					}
				}
				if lsn := sn.LogSequenceNumber(); lsn != test.result {
					t.Errorf("wrong log sequence number: %v, expected: %v", lsn, test.result)
				}
				if f == nil {
					t.Error("file should not be nil")
				}
			} else {
				if err != ErrNoSegmentFile {
					t.Error("should not find log sequence number")
				}
				if f != nil {
					t.Error("file should be nil")
				}
			}
		})
	}
}

func TestSegmentList_OpenNextSegment(t *testing.T) {
	t.Parallel()

	fsys := mock.NewFileSystem(fstest.MapFS{
		NewSegmentName(1, 10).String(): {},
		NewSegmentName(2, 20).String(): {},
	})
	sl := NewSegmentList(fsys, ".")

	if err := sl.Open(); err != nil {
		t.Fatal(err)
	}

	sn, f, err := sl.OpenNextSegment(30)
	if err != nil {
		t.Fatal(err)
	}

	if id := sn.ID(); id != 3 {
		t.Errorf("wrong id: %v, expected: %v", id, 3)
	}

	if f == nil {
		t.Error("file should not be nil")
	}

	if names := sl.Names(); len(names) <= 3 && !slices.Contains(names, sn) {
		t.Errorf("segment list does not contain segment: %v", sn)
	}
}

func TestSegmentList_RemoveSegmentBeforeLSN(t *testing.T) {
	t.Parallel()

	fsys := mock.NewFileSystem(fstest.MapFS{
		NewSegmentName(1, 10).String(): {},
		NewSegmentName(2, 20).String(): {},
	})
	sl := NewSegmentList(fsys, ".")

	if err := sl.Open(); err != nil {
		t.Fatal(err)
	}

	if err := sl.RemoveSegmentBeforeLSN(20); err != nil {
		t.Fatal(err)
	}

	if names := sl.Names(); !reflect.DeepEqual(names, []SegmentName{NewSegmentName(2, 20)}) {
		t.Errorf("segment list is incorrect: %v", names)
	}

	// Open again to verify
	if err := sl.Open(); err != nil {
		t.Fatal(err)
	}

	if names := sl.Names(); !reflect.DeepEqual(names, []SegmentName{NewSegmentName(2, 20)}) {
		t.Errorf("segment list is incorrect: %v", names)
	}
}

func createNames(start SegmentID, count SegmentID) []SegmentName {
	names := make([]SegmentName, 0, count)
	for i := range count {
		names = append(names, NewSegmentName(start+i, 0))
	}
	return names
}

func createFiles(start SegmentID, count SegmentID) *mock.FileSystem {
	fsys := &mock.FileSystem{MapFS: fstest.MapFS{}}
	for i := range count {
		fsys.MapFS[NewSegmentName(start+i, 0).String()] = &fstest.MapFile{}
	}
	return fsys
}
