package transaction

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/record"
)

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
