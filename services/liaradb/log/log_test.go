package log

import (
	"reflect"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/file/disk"
	"github.com/liaradb/liaradb/log/record"
	"github.com/liaradb/liaradb/log/segment"
)

func TestLog_EmptyReader(t *testing.T) {
	t.Parallel()

	fsys, dir := createFiles(t)

	l := NewLog(256, 2, fsys, dir)
	if err := l.Open(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}()

	for _, err := range l.Iterate(0) {
		if err != segment.ErrNoSegmentFile {
			t.Error("should have no files")
		}
	}
}

func TestLog_Iterate(t *testing.T) {
	t.Parallel()

	fsys, dir := createFiles(t)

	l := NewLog(256, 2, fsys, dir)
	if err := l.Open(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}()

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	records, _ := createRecords(100)
	var lsn record.LogSequenceNumber
	var err error
	for _, rec := range records {
		lsn, err = l.Append(rec)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err = l.Flush(lsn); err != nil {
		t.Fatal(err)
	}

	i := 0
	for rc, err := range l.Iterate(0) {
		if err != nil {
			t.Fatal(err)
		}

		rec := records[i]

		if !reflect.DeepEqual(rc, rec) {
			t.Error("records do not match")
		}
		i++
	}
	if i != 100 {
		t.Errorf("incorrect count: %v, expected: %v", i, 100)
	}
}

func TestLog_Recover(t *testing.T) {
	t.Parallel()

	fsys, dir := createFiles(t)
	records, _ := createRecords(2)

	t.Run("should append and flush", func(t *testing.T) {
		l := NewLog(256, 2, fsys, dir)
		if err := l.Open(); err != nil {
			t.Fatal(err)
		}

		if err := l.StartWriter(); err != nil {
			t.Fatal(err)
		}

		lsn1, err := l.Append(records[0])
		if err != nil {
			t.Fatal(err)
		}

		if err := l.Flush(lsn1); err != nil {
			t.Fatal(err)
		}

		lsn2, err := l.Append(records[1])
		if err != nil {
			t.Fatal(err)
		}

		if err := l.Flush(lsn2); err != nil {
			t.Fatal(err)
		}

		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should recover", func(t *testing.T) {
		l := NewLog(256, 2, fsys, dir)
		if err := l.Open(); err != nil {
			t.Fatal(err)
		}

		it, err := l.Recover()
		if err != nil {
			t.Fatal(err)
		}

		i := 0
		for rc := range it {
			if err != nil {
				t.Fatal(err)
			}

			rec := records[i]

			if !reflect.DeepEqual(rc, rec) {
				t.Error("records do not match")
			}
			i++
		}
		if i != 2 {
			t.Errorf("incorrect count: %v, expected: %v", i, 2)
		}
	})
}

func TestLog_RecoverMany(t *testing.T) {
	t.Parallel()

	fsys, dir := createFiles(t)

	var aCount1 record.LogSequenceNumber = 100
	var aCount2 record.LogSequenceNumber = 1
	aCount := aCount1 + aCount2
	records1, _ := createRecords(aCount1)
	records2, _ := createRecords(aCount2)
	records := append(records1, records2...)

	t.Run("should append and flush", func(t *testing.T) {
		l := NewLog(256, 2, fsys, dir)
		if err := l.Open(); err != nil {
			t.Fatal(err)
		}

		if err := l.StartWriter(); err != nil {
			t.Fatal(err)
		}

		var lsn record.LogSequenceNumber
		var err error
		for _, rec := range records1 {
			lsn, err = l.Append(rec)
			if err != nil {
				t.Fatal(err)
			}
		}

		if err = l.Flush(lsn); err != nil {
			t.Fatal(err)
		}

		i := 0
		for rc, err := range l.Iterate(0) {
			if err != nil {
				t.Fatal(err)
			}

			rec := records1[i]

			if !reflect.DeepEqual(rc, rec) {
				t.Error("records do not match")
			}
			i++
		}
		if i != int(aCount1) {
			t.Errorf("incorrect count: %v, expected: %v", i, aCount1)
		}

		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should append and flush more and iterate", func(t *testing.T) {
		t.Skip()
		l := NewLog(256, 2, fsys, dir)
		if err := l.Open(); err != nil {
			t.Fatal(err)
		}

		if err := l.StartWriter(); err != nil {
			t.Fatal(err)
		}

		var lsn record.LogSequenceNumber
		var err error
		for _, rec := range records2 {
			lsn, err = l.Append(rec)
			if err != nil {
				t.Fatal(err)
			}
		}

		if err := l.Flush(lsn); err != nil {
			t.Fatal(err)
		}

		i := 0
		for rc, err := range l.Iterate(0) {
			if err != nil {
				t.Fatal(err)
			}

			rec := records[i]

			if !reflect.DeepEqual(rc, rec) {
				t.Errorf("records do not match: %v, expected: %v",
					rc.LogSequenceNumber(),
					rec.LogSequenceNumber())
			}
			i++
		}
		if i != int(aCount) {
			t.Errorf("incorrect count: %v, expected: %v", i, aCount)
		}

		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestLog_Reverse(t *testing.T) {
	t.Parallel()

	fsys, dir := createFiles(t)

	l := NewLog(256, 2, fsys, dir)
	if err := l.Open(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := l.Close(); err != nil {
			t.Error(err)
		}
	}()

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	records, _ := createRecords(100)
	var lsn record.LogSequenceNumber
	var err error
	for _, rec := range records {
		lsn, err = l.Append(rec)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err = l.Flush(lsn); err != nil {
		t.Fatal(err)
	}

	slices.Reverse(records)
	i := 0
	for rc, err := range l.Reverse() {
		if err != nil {
			t.Fatal(err)
		}

		rec := records[i]

		if !reflect.DeepEqual(rc, rec) {
			t.Error("records do not match")
		}
		i++
	}
	if i != 100 {
		t.Errorf("incorrect count: %v, expected: %v", i, 100)
	}
}

func createFiles(t *testing.T) (file.FileSystem, string) {
	return &disk.FileSystem{}, t.TempDir()
	// return &mock.FileSystem{MapFS: fstest.MapFS{}}, "."
}
