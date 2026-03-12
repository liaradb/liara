package segment

import (
	"fmt"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/mock"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/testing/filetesting"
)

func TestList_Open(t *testing.T) {
	t.Parallel()

	var count SegmentID = 10

	t.Run("should list segments", func(t *testing.T) {
		t.Parallel()

		fsys := createFiles(t, 0, count)
		sl := NewList(fsys, dir)

		if err := sl.Open(); err != nil {
			t.Fatal(err)
		}

		names := sl.Names()

		want := createNames(0, count)
		if !slices.Equal(want, names) {
			t.Errorf("files do not match:\n\t%v,\nexpected:\n\t%v", names, want)
		}
	})

	t.Run("should list segments in order", func(t *testing.T) {
		t.Parallel()

		fsys := createFiles(t, 9998, count)
		sl := NewList(fsys, dir)

		if err := sl.Open(); err != nil {
			t.Fatal(err)
		}

		names := sl.Names()

		want := createNames(9998, count)
		if !slices.Equal(want, names) {
			t.Errorf("files do not match:\n\t%v,\nexpected:\n\t%v", names, want)
		}
	})
}

func TestList_OpenLatestSegment(t *testing.T) {
	t.Parallel()

	for message, test := range map[string]struct {
		result SegmentID
		fsys   file.FileSystem
	}{
		"should handle no files": {0, mock.NewFileSystem(fstest.MapFS{
			fmt.Sprintf("%v/", dir): &fstest.MapFile{},
		})},
		"should handle one file": {1, mock.NewFileSystem(fstest.MapFS{
			createPath(NewSegmentName(1, record.NewLogSequenceNumber(10))): {},
		})},
		"should handle multiple files": {2, mock.NewFileSystem(fstest.MapFS{
			createPath(NewSegmentName(1, record.NewLogSequenceNumber(10))): {},
			createPath(NewSegmentName(2, record.NewLogSequenceNumber(20))): {},
		})},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()

			sl := NewList(test.fsys, dir)

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
				t.Errorf("list does not contain segment: %v", sn)
			}
		})
	}

	t.Run("should close previous file", func(t *testing.T) {
		t.Parallel()

		fsys := mock.NewFileSystem(fstest.MapFS{
			createPath(NewSegmentName(1, record.NewLogSequenceNumber(10))): {},
			createPath(NewSegmentName(2, record.NewLogSequenceNumber(20))): {},
		})
		sl := NewList(fsys, dir)

		_, f, err := sl.OpenLatestSegment()
		if err != nil {
			t.Fatal(err)
		}

		if _, _, err = sl.OpenNextSegment(record.NewLogSequenceNumber(30)); err != nil {
			t.Fatal(err)
		}

		if _, _, err := sl.OpenLatestSegment(); err != nil {
			t.Fatal(err)
		}

		if f.(*mock.File).IsOpen() {
			t.Error("previous file should be closed")
		}
	})
}

func TestList_IterateFromLSN(t *testing.T) {
	t.Parallel()

	sn0 := NewSegmentName(1, record.NewLogSequenceNumber(10))
	sn1 := NewSegmentName(2, record.NewLogSequenceNumber(20))
	sn2 := NewSegmentName(3, record.NewLogSequenceNumber(30))
	names := []SegmentName{sn0, sn1, sn2}
	fsys := mock.NewFileSystem(fstest.MapFS{
		createPath(sn0): {},
		createPath(sn1): {},
		createPath(sn2): {},
	})
	sl := NewList(fsys, dir)

	c := 0
	n := make([]SegmentName, 0, 3)
	for f, err := range sl.IterateFromLSN(record.NewLogSequenceNumber(10)) {
		if err != nil {
			t.Fatal(err)
		}

		m, _ := f.(*mock.File).Stat()
		n = append(n, ParseSegmentName(m.Name()))
		c++
	}
	if c != 3 {
		t.Errorf("incorrect count: %v, expected: %v", c, 3)
	}
	if !slices.Equal(names, n) {
		t.Error("names do not match")
	}
}

func TestList_OpenSegmentForLSN(t *testing.T) {
	t.Parallel()

	fsys := mock.NewFileSystem(fstest.MapFS{
		createPath(NewSegmentName(1, record.NewLogSequenceNumber(10))): {},
		createPath(NewSegmentName(2, record.NewLogSequenceNumber(20))): {},
	})
	sl := NewList(fsys, dir)

	for message, test := range map[string]struct {
		search uint64
		found  bool
		result uint64
	}{
		"should not find low value": {1, false, 0},
		"should find exact value":   {10, true, 10},
		"should find middle value":  {15, true, 10},
		"should find high value":    {50, true, 20},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()

			sn, f, err := sl.OpenSegmentForLSN(record.NewLogSequenceNumber(test.search))
			if test.found {
				if err != nil {
					if err == ErrNoSegmentFile {
						t.Error("should find log sequence number")
					} else {
						t.Error(err)
					}
				}

				result := record.NewLogSequenceNumber(test.result)
				if lsn := sn.LogSequenceNumber(); lsn != result {
					t.Errorf("wrong log sequence number: %v, expected: %v", lsn, result)
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

	t.Run("should close previous file", func(t *testing.T) {
		t.Parallel()

		fsys := mock.NewFileSystem(fstest.MapFS{
			createPath(NewSegmentName(1, record.NewLogSequenceNumber(10))): {},
			createPath(NewSegmentName(2, record.NewLogSequenceNumber(20))): {},
		})
		sl := NewList(fsys, dir)

		_, f, err := sl.OpenSegmentForLSN(record.NewLogSequenceNumber(10))
		if err != nil {
			t.Fatal(err)
		}

		if _, _, err := sl.OpenSegmentForLSN(record.NewLogSequenceNumber(20)); err != nil {
			t.Fatal(err)
		}

		if f.(*mock.File).IsOpen() {
			t.Error("previous file should be closed")
		}
	})
}

func TestList_OpenNextSegment(t *testing.T) {
	t.Parallel()

	t.Run("should open next segment", func(t *testing.T) {
		fsys := mock.NewFileSystem(fstest.MapFS{
			createPath(NewSegmentName(1, record.NewLogSequenceNumber(10))): {},
			createPath(NewSegmentName(2, record.NewLogSequenceNumber(20))): {},
		})
		sl := NewList(fsys, dir)

		sn, f, err := sl.OpenNextSegment(record.NewLogSequenceNumber(30))
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
			t.Errorf("list does not contain segment: %v", sn)
		}

		// Open again to verify
		if err := sl.Open(); err != nil {
			t.Fatal(err)
		}

		if names := sl.Names(); len(names) <= 3 && !slices.Contains(names, sn) {
			t.Errorf("list does not contain segment: %v", sn)
		}
	})

	t.Run("should close previous file", func(t *testing.T) {
		t.Parallel()

		fsys := mock.NewFileSystem(fstest.MapFS{
			createPath(NewSegmentName(1, record.NewLogSequenceNumber(10))): {},
		})
		sl := NewList(fsys, dir)

		_, f, err := sl.OpenNextSegment(record.NewLogSequenceNumber(20))
		if err != nil {
			t.Fatal(err)
		}

		if _, _, err := sl.OpenNextSegment(record.NewLogSequenceNumber(30)); err != nil {
			t.Fatal(err)
		}

		if f.(*mock.File).IsOpen() {
			t.Error("previous file should be closed")
		}
	})
}

func TestList_OpenSegmentBeforeLSN(t *testing.T) {
	t.Parallel()

	t.Run("should open segment", func(t *testing.T) {
		fsys := mock.NewFileSystem(fstest.MapFS{
			createPath(NewSegmentName(1, record.NewLogSequenceNumber(10))): {},
			createPath(NewSegmentName(2, record.NewLogSequenceNumber(20))): {},
			createPath(NewSegmentName(3, record.NewLogSequenceNumber(30))): {},
		})
		sl := NewList(fsys, dir)

		sn, f, err := sl.OpenSegmentBeforeLSN(record.NewLogSequenceNumber(30))
		if err != nil {
			t.Fatal(err)
		}

		if id := sn.ID(); id != 2 {
			t.Errorf("wrong id: %v, expected: %v", id, 2)
		}

		if f == nil {
			t.Error("file should not be nil")
		}

		if names := sl.Names(); !slices.Contains(names, sn) {
			t.Errorf("list does not contain segment: %v", sn)
		}
	})

	t.Run("should close previous file", func(t *testing.T) {
		t.Parallel()

		fsys := mock.NewFileSystem(fstest.MapFS{
			createPath(NewSegmentName(1, record.NewLogSequenceNumber(10))): {},
			createPath(NewSegmentName(2, record.NewLogSequenceNumber(20))): {},
			createPath(NewSegmentName(3, record.NewLogSequenceNumber(30))): {},
		})
		sl := NewList(fsys, dir)

		_, f, err := sl.OpenSegmentBeforeLSN(record.NewLogSequenceNumber(30))
		if err != nil {
			t.Fatal(err)
		}

		if _, _, err := sl.OpenSegmentBeforeLSN(record.NewLogSequenceNumber(20)); err != nil {
			t.Fatal(err)
		}

		if f.(*mock.File).IsOpen() {
			t.Error("previous file should be closed")
		}
	})
}

func TestList_RemoveSegmentBeforeLSN(t *testing.T) {
	t.Parallel()

	fsys := mock.NewFileSystem(fstest.MapFS{
		createPath(NewSegmentName(1, record.NewLogSequenceNumber(10))): {},
		createPath(NewSegmentName(2, record.NewLogSequenceNumber(20))): {},
	})
	sl := NewList(fsys, dir)

	if err := sl.RemoveSegmentBeforeLSN(record.NewLogSequenceNumber(20)); err != nil {
		t.Fatal(err)
	}

	if names := sl.Names(); !slices.Equal(names, []SegmentName{NewSegmentName(2, record.NewLogSequenceNumber(20))}) {
		t.Errorf("list is incorrect: %v", names)
	}

	// Open again to verify
	if err := sl.Open(); err != nil {
		t.Fatal(err)
	}

	if names := sl.Names(); !slices.Equal(names, []SegmentName{NewSegmentName(2, record.NewLogSequenceNumber(20))}) {
		t.Errorf("list is incorrect: %v", names)
	}
}

func TestList_Reverse(t *testing.T) {
	t.Parallel()

	sn0 := NewSegmentName(1, record.NewLogSequenceNumber(10))
	sn1 := NewSegmentName(2, record.NewLogSequenceNumber(20))
	sn2 := NewSegmentName(3, record.NewLogSequenceNumber(30))
	names := []SegmentName{sn0, sn1, sn2}
	slices.Reverse(names)
	fsys := mock.NewFileSystem(fstest.MapFS{
		createPath(sn0): {},
		createPath(sn1): {},
		createPath(sn2): {},
	})
	sl := NewList(fsys, dir)

	c := 0
	n := make([]SegmentName, 0, 3)
	for f, err := range sl.Reverse() {
		if err != nil {
			t.Fatal(err)
		}

		m, _ := f.(*mock.File).Stat()
		n = append(n, ParseSegmentName(m.Name()))
		c++
	}
	if c != 3 {
		t.Errorf("incorrect count: %v, expected: %v", c, 3)
	}
	if !slices.Equal(names, n) {
		t.Error("names do not match")
	}
}

func createNames(start SegmentID, count SegmentID) []SegmentName {
	names := make([]SegmentName, 0, count)
	for i := range count {
		names = append(names, NewSegmentName(start+i, record.NewLogSequenceNumber(0)))
	}
	return names
}

func createFiles(t *testing.T, start SegmentID, count SegmentID) *mock.FileSystem {
	fsys := filetesting.NewMockFileSystem(t, fstest.MapFS{})
	for i := range count {
		fsys.MapFS[createPath(NewSegmentName(start+i, record.NewLogSequenceNumber(0)))] = &fstest.MapFile{}
	}
	return fsys
}

const dir = "log"

func createPath(sn SegmentName) string {
	return fmt.Sprintf("%v/%v", dir, sn)
}
