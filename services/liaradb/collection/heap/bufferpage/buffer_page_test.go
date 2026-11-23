package bufferpage

import (
	"testing"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/should"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
	"github.com/liaradb/liaradb/util/testing/txtesting"
)

func TestBufferPage(t *testing.T) {
	storagetesting.SyncTest(t, 100, 1024, func(t *testing.T, s storagetesting.Storage) {
		tx := txtesting.New(s.Storage)
		defer tx.Release()

		fbid := link.NewFileName("test").FileBlockID(link.BlockID(0))
		b, err := tx.Pin(t.Context(), fbid)
		if err != nil {
			t.Fatal(err)
		}

		bp := New(b, tx)
		defer bp.Release()

		if i := bp.FileBlockID(); i != fbid {
			t.Errorf("incorrect file block id: %v, expected: %v", i, fbid)
		}

		if lsn := bp.LogSequenceNumber(); lsn != logpage.NewLogSequenceNumber(0) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(0))
		}

		bp.SetLogSequenceNumber(logpage.NewLogSequenceNumber(1))

		if lsn := bp.LogSequenceNumber(); lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		if c := bp.Count(); c != 0 {
			t.Errorf("incorrect count: %v, expected: %v", c, 0)
		}

		{
			h, d := bp.Next(8)
			h[0] = 1
			d[0] = 2

			if b.Dirty() {
				t.Error("should not be dirty")
			}

			bp.CommitAppend(8)

			if !b.Dirty() {
				t.Error("should be dirty")
			}
		}

		if c := bp.Count(); c != 1 {
			t.Errorf("incorrect count: %v, expected: %v", c, 1)
		}

		{
			h, d := bp.Slot(0)
			if h[0] != 1 {
				t.Errorf("incorrect data: %v, expected: %v", h[0], 1)
			}
			if d[0] != 2 {
				t.Errorf("incorrect data: %v, expected: %v", d[0], 2)
			}
		}

		{
			bp.Delete(0)
			should.Panic(t, "slot should not exist", func() {
				_, _ = bp.Slot(0)
			})
		}
	})
}

func TestBufferPage_CommitUpdate(t *testing.T) {
	storagetesting.SyncTest(t, 100, 1024, func(t *testing.T, s storagetesting.Storage) {
		tx := txtesting.New(s.Storage)
		defer tx.Release()

		fbid := link.NewFileName("test").FileBlockID(link.BlockID(0))
		b0, err := tx.Pin(t.Context(), fbid)
		if err != nil {
			t.Fatal(err)
		}

		bp0 := New(b0, tx)
		defer bp0.Release()

		if b0.Dirty() {
			t.Error("should not be dirty")
		}

		bp0.CommitUpdate()

		if !b0.Dirty() {
			t.Error("should be dirty")
		}
	})
}

func TestBufferPage_Shadow(t *testing.T) {
	storagetesting.SyncTest(t, 100, 1024, func(t *testing.T, s storagetesting.Storage) {
		tx := txtesting.New(s.Storage)
		defer tx.Release()

		fbid := link.NewFileName("test").FileBlockID(link.BlockID(0))
		b0, err := tx.Pin(t.Context(), fbid)
		if err != nil {
			t.Fatal(err)
		}

		bp0 := New(b0, tx)
		defer bp0.Release()

		b1, err := tx.Pin(t.Context(), fbid)
		if err != nil {
			t.Fatal(err)
		}

		bp1 := New(b1, tx)
		defer bp1.Release()

		{
			bp0.SetLogSequenceNumber(logpage.NewLogSequenceNumber(1))
			h, d := bp0.Next(8)
			h[0] = 1
			d[0] = 2
			bp0.CommitAppend(8)
		}

		bp1.Shadow(bp0)
		if c := bp1.Count(); c != 1 {
			t.Errorf("incorrect count: %v, expected: %v", c, 1)
		}

		if lsn := bp1.LogSequenceNumber(); lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		{
			h, d := bp1.Slot(0)
			if h[0] != 1 {
				t.Errorf("incorrect data: %v, expected: %v", h[0], 1)
			}
			if d[0] != 2 {
				t.Errorf("incorrect data: %v, expected: %v", d[0], 2)
			}
		}
	})
}
