package transaction

import (
	"fmt"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/record"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

// TODO: Ensure Transaction pins BufferList
// TODO: Verify Transaction.Release
func TestTransaction_Run(t *testing.T) {
	t.Parallel()

	t.Run("should commit", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			m, _, iterator := createManagerIterator(t)
			ctx := t.Context()

			tx, err := m.Next(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if err := tx.Run(t.Context(), func() error {
				_, err := tx.Insert(t.Context(), link.FileRecordID{}, nil)
				return err
			}); err != nil {
				t.Fatal(err)
			}

			want := []record.Action{
				record.ActionStart,
				record.ActionInsert,
				record.ActionCommit,
			}

			actions := make([]record.Action, 0, len(want))
			for r := range iterator() {
				actions = append(actions, r.Action())
			}

			if !slices.Equal(actions, want) {
				t.Errorf("incorrect actions: %v, expected: %v", actions, want)
			}
		})
	})

	t.Run("should rollback", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			m, _, iterator := createManagerIterator(t)
			ctx := t.Context()

			tx, err := m.Next(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if err := tx.Run(t.Context(), func() error {
				if _, err := tx.Insert(t.Context(), link.FileRecordID{}, nil); err != nil {
					t.Fatal(err)
				}
				return fmt.Errorf("error")
			}); err == nil {
				t.Error("should return error")
			}

			want := []record.Action{
				record.ActionStart,
				record.ActionInsert,
				record.ActionRollback,
			}

			actions := make([]record.Action, 0, len(want))
			for r := range iterator() {
				actions = append(actions, r.Action())
			}

			if !slices.Equal(actions, want) {
				t.Errorf("incorrect actions: %v, expected: %v", actions, want)
			}
		})
	})
}

func TestTransaction_RunResult(t *testing.T) {
	t.Parallel()

	t.Run("should commit", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			m, _, iterator := createManagerIterator(t)
			ctx := t.Context()

			tx, err := m.Next(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if n, err := tx.RunResult(t.Context(), func() (int, error) {
				_, err := tx.Insert(t.Context(), link.FileRecordID{}, nil)
				return 1, err
			}); err != nil {
				t.Fatal(err)
			} else if n != 1 {
				t.Errorf("incorrect result: %v, expected: %v", n, 1)
			}

			want := []record.Action{
				record.ActionStart,
				record.ActionInsert,
				record.ActionCommit,
			}

			actions := make([]record.Action, 0, len(want))
			for r := range iterator() {
				actions = append(actions, r.Action())
			}

			if !slices.Equal(actions, want) {
				t.Errorf("incorrect actions: %v, expected: %v", actions, want)
			}
		})
	})

	t.Run("should rollback", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			m, _, iterator := createManagerIterator(t)
			ctx := t.Context()

			tx, err := m.Next(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if n, err := tx.RunResult(t.Context(), func() (int, error) {
				if _, err := tx.Insert(t.Context(), link.FileRecordID{}, nil); err != nil {
					t.Fatal(err)
				}
				return 1, fmt.Errorf("error")
			}); err == nil {
				t.Error("should return error")
			} else if n != 1 {
				t.Errorf("incorrect result: %v, expected: %v", n, 1)
			}

			want := []record.Action{
				record.ActionStart,
				record.ActionInsert,
				record.ActionRollback,
			}

			actions := make([]record.Action, 0, len(want))
			for r := range iterator() {
				actions = append(actions, r.Action())
			}

			if !slices.Equal(actions, want) {
				t.Errorf("incorrect actions: %v, expected: %v", actions, want)
			}
		})
	})
}

func TestTransaction_Append(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		m, _, iterator := createManagerIterator(t)
		ctx := t.Context()

		tx, err := m.Next(ctx)
		if err != nil {
			t.Fatal(err)
		}

		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		if lsn, err := tx.Append(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data); err != nil {
			t.Error(err)
		} else if lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		count := 0
		for r := range iterator() {
			if a, want := r.Action(), record.ActionAppend; a != want {
				t.Errorf("incorrect : %v, expected: %v", a, want)
			}

			count++
		}

		if count != 1 {
			t.Errorf("incorrect count: %v, expected: %v", count, 1)
		}
	})
}

func TestTransaction_Insert(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		m, _, iterator := createManagerIterator(t)
		ctx := t.Context()

		tx, err := m.Next(ctx)
		if err != nil {
			t.Fatal(err)
		}

		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		if lsn, err := tx.Insert(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data); err != nil {
			t.Error(err)
		} else if lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		count := 0
		for r := range iterator() {
			if a, want := r.Action(), record.ActionInsert; a != want {
				t.Errorf("incorrect : %v, expected: %v", a, want)
			}

			count++
		}

		if count != 1 {
			t.Errorf("incorrect count: %v, expected: %v", count, 1)
		}
	})
}

func TestTransaction_Update(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		m, _, iterator := createManagerIterator(t)
		ctx := t.Context()

		tx, err := m.Next(ctx)
		if err != nil {
			t.Fatal(err)
		}

		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		prev := []byte{9, 10, 11, 12, 13, 14, 15, 16}
		if lsn, err := tx.Update(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data, prev); err != nil {
			t.Error(err)
		} else if lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		count := 0
		for r := range iterator() {
			if a, want := r.Action(), record.ActionUpdate; a != want {
				t.Errorf("incorrect : %v, expected: %v", a, want)
			}

			count++
		}

		if count != 1 {
			t.Errorf("incorrect count: %v, expected: %v", count, 1)
		}
	})
}

func TestTransaction_Clear(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		m, _, iterator := createManagerIterator(t)
		ctx := t.Context()

		tx, err := m.Next(ctx)
		if err != nil {
			t.Fatal(err)
		}

		prev := []byte{9, 10, 11, 12, 13, 14, 15, 16}
		if lsn, err := tx.Clear(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), prev); err != nil {
			t.Error(err)
		} else if lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		count := 0
		for r := range iterator() {
			if a, want := r.Action(), record.ActionClear; a != want {
				t.Errorf("incorrect : %v, expected: %v", a, want)
			}

			count++
		}

		if count != 1 {
			t.Errorf("incorrect count: %v, expected: %v", count, 1)
		}
	})
}

func TestTransaction_UpdateHeader(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		m, _, iterator := createManagerIterator(t)
		ctx := t.Context()

		tx, err := m.Next(ctx)
		if err != nil {
			t.Fatal(err)
		}

		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		prev := []byte{9, 10, 11, 12, 13, 14, 15, 16}
		if lsn, err := tx.UpdateHeader(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data, prev); err != nil {
			t.Error(err)
		} else if lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		count := 0
		for r := range iterator() {
			if a, want := r.Action(), record.ActionUpdateHeader; a != want {
				t.Errorf("incorrect : %v, expected: %v", a, want)
			}

			count++
		}

		if count != 1 {
			t.Errorf("incorrect count: %v, expected: %v", count, 1)
		}
	})
}

func TestTransaction_Pin(t *testing.T) {
	storagetesting.SyncTest(t, 200, 1024, func(t *testing.T, s storagetesting.Storage) {
		m, _ := createManager(t)
		tx, err := m.Next(t.Context())
		if err != nil {
			t.Fatal(err)
		}

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
