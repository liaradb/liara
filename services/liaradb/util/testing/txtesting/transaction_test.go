package txtesting

import (
	"fmt"
	"testing"

	"github.com/liaradb/liaradb/domain/infrastructure/tenantfolder"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/record"
)

func TestTransaction__Default(t *testing.T) {
	t.Parallel()

	tx := Transaction{}
	if l, want := tx.LogSequenceNumber(), logpage.NewLogSequenceNumber(0); l != want {
		t.Errorf("incorrect log sequence number: %v, expected: %v", l, want)
	}
	if l, want := len(tx.Records()), 0; l != want {
		t.Errorf("incorrect records: %v, expected: %v", l, want)
	}
}

func Test__Lock(t *testing.T) {
	t.Parallel()

	t.Run("should not be locked by default", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}

		tf := tenantfolder.NewTenantFolder(value.NewTenantID())
		tn := tf.Tenant()

		if tx.IsLocked(tn, "key") {
			t.Error("should not be locked")
		}

		if tx.IsSLocked(tn, "key") {
			t.Error("should not be s locked")
		}

		if tx.IsXLocked(tn, "key") {
			t.Error("should not be x locked")
		}
	})

	t.Run("should handle multiple table locks", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}

		tf := tenantfolder.NewTenantFolder(value.NewTenantID())
		tn0 := tf.Tenant()

		if err := tx.SLock(t.Context(), tn0, "key"); err != nil {
			t.Error(err)
		}

		tn1 := tf.Outbox()

		if tx.IsLocked(tn1, "key") {
			t.Error("should not be locked")
		}

		if tx.IsSLocked(tn1, "key") {
			t.Error("should not be s locked")
		}

		if tx.IsXLocked(tn1, "key") {
			t.Error("should not be x locked")
		}
	})

	t.Run("should s lock", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}

		tf := tenantfolder.NewTenantFolder(value.NewTenantID())
		tn := tf.Tenant()

		if err := tx.SLock(t.Context(), tn, "key"); err != nil {
			t.Error(err)
		}

		if !tx.IsLocked(tn, "key") {
			t.Error("should be locked")
		}

		if !tx.IsSLocked(tn, "key") {
			t.Error("should be s locked")
		}

		if tx.IsXLocked(tn, "key") {
			t.Error("should not be x locked")
		}
	})

	t.Run("should x lock", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}

		tf := tenantfolder.NewTenantFolder(value.NewTenantID())
		tn := tf.Tenant()

		if err := tx.SLock(t.Context(), tn, "key"); err != nil {
			t.Error(err)
		}

		if err := tx.XLock(t.Context(), tn, "key"); err != nil {
			t.Error(err)
		}

		if !tx.IsLocked(tn, "key") {
			t.Error("should be locked")
		}

		if tx.IsSLocked(tn, "key") {
			t.Error("should not be s locked")
		}

		if !tx.IsXLocked(tn, "key") {
			t.Error("should be x locked")
		}
	})

	t.Run("should upgrade to x lock", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}

		tf := tenantfolder.NewTenantFolder(value.NewTenantID())
		tn := tf.Tenant()

		if err := tx.SLock(t.Context(), tn, "key"); err != nil {
			t.Error(err)
		}

		if err := tx.XLock(t.Context(), tn, "key"); err != nil {
			t.Error(err)
		}

		if !tx.IsLocked(tn, "key") {
			t.Error("should be locked")
		}

		if tx.IsSLocked(tn, "key") {
			t.Error("should not be s locked")
		}

		if !tx.IsXLocked(tn, "key") {
			t.Error("should be x locked")
		}
	})

	t.Run("should not downgrade to s lock", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}

		tf := tenantfolder.NewTenantFolder(value.NewTenantID())
		tn := tf.Tenant()

		if err := tx.XLock(t.Context(), tn, "key"); err != nil {
			t.Error(err)
		}

		if err := tx.SLock(t.Context(), tn, "key"); err != nil {
			t.Error(err)
		}

		if !tx.IsLocked(tn, "key") {
			t.Error("should be locked")
		}

		if tx.IsSLocked(tn, "key") {
			t.Error("should not be s locked")
		}

		if !tx.IsXLocked(tn, "key") {
			t.Error("should be x locked")
		}
	})
}

func TestTransaction_Append(t *testing.T) {
	t.Parallel()

	t.Run("should append", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		if lsn, err := tx.Append(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data); err != nil {
			t.Error(err)
		} else if lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		if l, want := len(tx.Records()), 1; l != want {
			t.Errorf("incorrect length: %v, expected: %v", l, want)
		}

		r := tx.Records()[0]
		if a, want := r.Action(), record.ActionAppend; a != want {
			t.Errorf("incorrect : %v, expected: %v", a, want)
		}
	})

	t.Run("should return error if error is set", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}
		tx.SetError(fmt.Errorf("error"))

		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		if _, err := tx.Append(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data); err == nil {
			t.Error("should return error")
		}
	})
}

func TestTransaction_Insert(t *testing.T) {
	t.Parallel()

	t.Run("should insert", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		if lsn, err := tx.Insert(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data); err != nil {
			t.Error(err)
		} else if lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		if l, want := len(tx.Records()), 1; l != want {
			t.Errorf("incorrect length: %v, expected: %v", l, want)
		}

		r := tx.Records()[0]
		if a, want := r.Action(), record.ActionInsert; a != want {
			t.Errorf("incorrect : %v, expected: %v", a, want)
		}
	})

	t.Run("should return error if error is set", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}
		tx.SetError(fmt.Errorf("error"))

		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		if _, err := tx.Insert(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data); err == nil {
			t.Error("should return error")
		}
	})
}

func TestTransaction_Update(t *testing.T) {
	t.Parallel()

	t.Run("should insert", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		prev := []byte{9, 10, 11, 12, 13, 14, 15, 16}
		if lsn, err := tx.Update(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data, prev); err != nil {
			t.Error(err)
		} else if lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		if l, want := len(tx.Records()), 1; l != want {
			t.Errorf("incorrect length: %v, expected: %v", l, want)
		}

		r := tx.Records()[0]
		if a, want := r.Action(), record.ActionUpdate; a != want {
			t.Errorf("incorrect : %v, expected: %v", a, want)
		}
	})

	t.Run("should return error if error is set", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}
		tx.SetError(fmt.Errorf("error"))

		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		prev := []byte{9, 10, 11, 12, 13, 14, 15, 16}
		if _, err := tx.Update(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data, prev); err == nil {
			t.Error("should return error")
		}
	})
}

func TestTransaction_Clear(t *testing.T) {
	t.Parallel()

	t.Run("should insert", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}
		prev := []byte{9, 10, 11, 12, 13, 14, 15, 16}
		if lsn, err := tx.Clear(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), prev); err != nil {
			t.Error(err)
		} else if lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		if l, want := len(tx.Records()), 1; l != want {
			t.Errorf("incorrect length: %v, expected: %v", l, want)
		}

		r := tx.Records()[0]
		if a, want := r.Action(), record.ActionClear; a != want {
			t.Errorf("incorrect : %v, expected: %v", a, want)
		}
	})

	t.Run("should return error if error is set", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}
		tx.SetError(fmt.Errorf("error"))

		prev := []byte{9, 10, 11, 12, 13, 14, 15, 16}
		if _, err := tx.Clear(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), prev); err == nil {
			t.Error("should return error")
		}
	})
}

func TestTransaction_UpdateHeader(t *testing.T) {
	t.Parallel()

	t.Run("should insert", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		prev := []byte{9, 10, 11, 12, 13, 14, 15, 16}
		if lsn, err := tx.UpdateHeader(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data, prev); err != nil {
			t.Error(err)
		} else if lsn != logpage.NewLogSequenceNumber(1) {
			t.Errorf("incorrect log sequence number: %v, expected: %v", lsn, logpage.NewLogSequenceNumber(1))
		}

		if l, want := len(tx.Records()), 1; l != want {
			t.Errorf("incorrect length: %v, expected: %v", l, want)
		}

		r := tx.Records()[0]
		if a, want := r.Action(), record.ActionUpdateHeader; a != want {
			t.Errorf("incorrect : %v, expected: %v", a, want)
		}
	})

	t.Run("should return error if error is set", func(t *testing.T) {
		t.Parallel()

		tx := Transaction{}
		tx.SetError(fmt.Errorf("error"))

		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		prev := []byte{9, 10, 11, 12, 13, 14, 15, 16}
		if _, err := tx.UpdateHeader(t.Context(), link.NewFileRecordID(link.NewFileBlockID(link.NewFileName("file"), link.BlockID(1)), link.SlotID(2)), data, prev); err == nil {
			t.Error("should return error")
		}
	})
}
