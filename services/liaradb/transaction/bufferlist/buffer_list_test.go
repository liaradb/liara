package bufferlist

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestBufferList(t *testing.T) {
	storagetesting.SyncTest(t, 2, 1024, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		bl := NewBufferList(s.Storage)

		fn := link.NewFileName("testfile")

		b0, err := bl.Pin(ctx, fn.FileBlockID(1))
		if err != nil {
			t.Fatal(err)
		}
		if i := b0.FileBlockID().BlockID(); i != 1 {
			t.Errorf("incorrect result: expected %v, recieved %v", 1, i)
		}
		if p := b0.Pins(); p != 1 {
			t.Errorf("incorrect pins: %v, expected: %v", p, 1)
		}

		if r := b0.Reads(); r != 1 {
			t.Errorf("incorrect reads: %v, expected: %v", r, 1)
		}

		b1, err := bl.Pin(ctx, fn.FileBlockID(2))
		if err != nil {
			t.Fatal(err)
		}
		if i := b1.FileBlockID().BlockID(); i != 2 {
			t.Fatalf("incorrect result: expected %v, recieved %v", 2, i)
		}
		if p := b1.Pins(); p != 1 {
			t.Errorf("incorrect pins: %v, expected: %v", p, 1)
		}
		if r := b1.Reads(); r != 1 {
			t.Errorf("incorrect reads: %v, expected: %v", r, 1)
		}

		bl.Release()

		synctest.Wait()
		if p := b0.Pins(); p != 0 {
			t.Errorf("incorrect pins: %v, expected: %v", p, 0)
		}
		if r := b0.Reads(); r != 1 {
			t.Errorf("incorrect reads: %v, expected: %v", r, 1)
		}
		if p := b1.Pins(); p != 0 {
			t.Errorf("incorrect pins: %v, expected: %v", p, 0)
		}
		if r := b1.Reads(); r != 1 {
			t.Errorf("incorrect reads: %v, expected: %v", r, 1)
		}
	})
}

func TestBufferList_Pin(t *testing.T) {
	storagetesting.SyncTest(t, 200, 1024, func(t *testing.T, s storagetesting.Storage) {
		tx := NewBufferList(s.Storage)
		defer tx.Release()

		fn := link.NewFileName("test")
		fbid0 := link.NewFileBlockID(fn, link.BlockID(0))
		b0, err := tx.Pin(t.Context(), fbid0)
		if err != nil {
			t.Fatal(err)
		}

		if v := b0.FileBlockID(); v != fbid0 {
			t.Errorf("incorrect block: %v, expected: %v", v, fbid0)
		}

		c0, err := tx.PinCurrent(t.Context(), fn)
		if err != nil {
			t.Fatal(err)
		}

		if v := c0.FileBlockID(); v != fbid0 {
			t.Errorf("incorrect block: %v, expected: %v", v, fbid0)
		}

		n0, err := tx.PinNext(t.Context(), fn)
		if err != nil {
			t.Fatal(err)
		}

		fbid1 := link.NewFileBlockID("test", link.BlockID(1))
		if v := n0.FileBlockID(); v != fbid1 {
			t.Errorf("incorrect block: %v, expected: %v", v, fbid1)
		}

		b1, err := tx.Pin(t.Context(), fbid1)
		if err != nil {
			t.Fatal(err)
		}

		if v := b1.FileBlockID(); v != fbid1 {
			t.Errorf("incorrect block: %v, expected: %v", v, fbid1)
		}

		c1, err := tx.PinCurrent(t.Context(), fn)
		if err != nil {
			t.Fatal(err)
		}

		if v := c1.FileBlockID(); v != fbid1 {
			t.Errorf("incorrect block: %v, expected: %v", v, fbid1)
		}
	})
}
