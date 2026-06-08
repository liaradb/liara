package pageiterator

import (
	"testing"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/recovery/pagequeue"
	"github.com/liaradb/liaradb/recovery/pagestorage"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/util/testing/filetesting"
)

func TestPageIterator(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	txid := record.NewTransactionID(2)
	now := record.NewTime(time.UnixMicro(1234567890))
	action := record.ActionInsert
	collection := record.CollectionEvent
	data := []byte("abcdef")
	reverse := []byte("fghij")

	t.Run("should run", func(t *testing.T) {
		t.Parallel()

		fsys := filetesting.New(nil)
		dir := "dir"

		const (
			size       = 128
			headerSize = 4
			slotSize   = 4
		)

		sl := segment.NewList(fsys, dir, size, 1)
		ps := pagestorage.New(sl)
		if err := ps.Init(); err != nil {
			t.Error(err)
		}

		pq := pagequeue.New(ps, size, headerSize, slotSize)

		numberOfRecords := 100

		for i := range numberOfRecords {
			rc := record.New(record.NewLogSequenceNumber(uint64(i)), tid, txid, now, action, collection, data, reverse)
			if err := pq.Append(rc); err != nil {
				t.Error(err)
			}
		}

		if err := pq.Flush(); err != nil {
			t.Error(err)
		}

		if err := sl.Close(); err != nil {
			t.Error(err)
		}

		sl = segment.NewList(fsys, dir, size, 1)
		pi := New(sl, size, headerSize, slotSize)

		c := 0
		for rc, err := range pi.Forward(record.NewLogSequenceNumber(0)) {
			if err != nil {
				t.Error(err)
			}

			want := uint64(c)
			if lsn := rc.LogSequenceNumber().Value(); lsn != want {
				t.Errorf("incorrect lsn: %v, expected: %v", lsn, want)
			}

			c++
		}

		if c != numberOfRecords {
			t.Errorf("incorrect count: %v, expected: %v", c, numberOfRecords)
		}
	})

	t.Run("should reverse", func(t *testing.T) {
		t.Parallel()

		fsys := filetesting.New(nil)
		dir := "dir"

		const (
			size       = 128
			headerSize = 4
			slotSize   = 4
		)

		sl := segment.NewList(fsys, dir, size, 1)
		ps := pagestorage.New(sl)
		if err := ps.Init(); err != nil {
			t.Error(err)
		}

		pq := pagequeue.New(ps, size, headerSize, slotSize)

		numberOfRecords := 100

		for i := range numberOfRecords {
			rc := record.New(record.NewLogSequenceNumber(uint64(i)), tid, txid, now, action, collection, data, reverse)
			if err := pq.Append(rc); err != nil {
				t.Error(err)
			}
		}

		if err := pq.Flush(); err != nil {
			t.Error(err)
		}

		if err := sl.Close(); err != nil {
			t.Error(err)
		}

		sl = segment.NewList(fsys, dir, size, 1)
		pi := New(sl, size, headerSize, slotSize)

		c := 0
		reverseIndex := numberOfRecords - 1
		for rc, err := range pi.Reverse() {
			if err != nil {
				t.Fatal(err)
			}

			want := uint64(reverseIndex)
			if lsn := rc.LogSequenceNumber().Value(); lsn != want {
				t.Errorf("incorrect lsn: %v, expected: %v", lsn, want)
			}

			c++
			reverseIndex--
		}

		if c != numberOfRecords {
			t.Errorf("incorrect count: %v, expected: %v", c, numberOfRecords)
		}
	})
}
