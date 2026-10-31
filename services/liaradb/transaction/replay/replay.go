package replay

import (
	"context"
	"iter"

	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/transaction/record"
	"github.com/liaradb/liaradb/util/set"
)

type Replay struct {
	s     *storage.Storage
	log   *log.Log
	txIDs set.Set[record.TransactionID]
}

func NewReplay(
	s *storage.Storage,
	log *log.Log,
) *Replay {
	return &Replay{
		s:   s,
		log: log,
	}
}

func (re *Replay) Recover(ctx context.Context) error {
	it, err := re.log.Recover()
	if err != nil {
		return err
	}

	for r := range it {
		if err := re.recoverRecord(ctx, r); err != nil {
			return err
		}
	}

	return nil
}

func (re *Replay) recoverRecord(ctx context.Context, r *record.Record) error {
	switch r.Action() {
	case record.ActionCheckpoint:
		return re.recoverCheckpoint(r)
	case record.ActionStart:
		return re.recoverStart(r)
	case record.ActionCommit:
		return re.recoverCommit()
	case record.ActionAppend:
		return re.recoverAppend(ctx, r)
	case record.ActionInsert:
		return re.recoverInsert(ctx, r)
	case record.ActionUpdateHeader:
		return re.recoverUpdateHeader(ctx, r)
	case record.ActionRemove:
		return re.recoverRemove(ctx, r)
	case record.ActionClear:
		return re.recordClear(ctx, r)
	case record.ActionRollback:
		return re.recoverRollback()
	case record.ActionUpdate:
		return re.recoverUpdate(ctx, r)
	default:
		return ErrActionUnknown
	}
}

func (re *Replay) recoverCheckpoint(r *record.Record) error {
	for txID := range re.txIDsFromData(r.Data()) {
		re.txIDs.Add(txID)
	}
	return nil
}

func (*Replay) txIDsFromData(
	data []byte,
) iter.Seq[record.TransactionID] {
	return func(yield func(record.TransactionID) bool) {
		for i := range len(data) % record.TransactionIDSize {
			start := i * record.TransactionIDSize
			item := data[start : start+record.TransactionIDSize]
			txID := record.TransactionID{}
			_ = txID.ReadData(item)
			if !yield(txID) {
				return
			}
		}
	}
}

func (re *Replay) recoverStart(r *record.Record) error {
	re.txIDs.Remove(r.TransactionID())
	return nil
}

func (re *Replay) recoverCommit() error {
	return nil
}

func (re *Replay) recoverAppend(ctx context.Context, r *record.Record) error {
	b, err := re.s.Request(ctx, r.FileRecordID().FileBlockID())
	if err != nil {
		return err
	}

	defer b.Release()

	p := re.getPage(r, b)
	if re.skipLSN(r, p) {
		return nil
	}

	p.RecoverAppend(r.Data())
	return nil
}

func (re *Replay) recoverInsert(ctx context.Context, r *record.Record) error {
	b, err := re.s.Request(ctx, r.FileRecordID().FileBlockID())
	if err != nil {
		return err
	}

	defer b.Release()

	p := re.getPage(r, b)
	if re.skipLSN(r, p) {
		return nil
	}

	p.RecoverInsert(r.FileRecordID().SlotID(), r.Data())
	return nil
}

func (re *Replay) recoverUpdateHeader(ctx context.Context, r *record.Record) error {
	b, err := re.s.Request(ctx, r.FileRecordID().FileBlockID())
	if err != nil {
		return err
	}

	defer b.Release()

	p := re.getPage(r, b)
	if re.skipLSN(r, p) {
		return nil
	}

	p.RecoverHeader(r.Data())
	return nil
}

func (re *Replay) recoverRemove(ctx context.Context, r *record.Record) error {
	b, err := re.s.Request(ctx, r.FileRecordID().FileBlockID())
	if err != nil {
		return err
	}

	defer b.Release()

	p := re.getPage(r, b)
	if re.skipLSN(r, p) {
		return nil
	}

	// TODO: Run remove method
	return nil
}

func (re *Replay) recoverUpdate(ctx context.Context, r *record.Record) error {
	b, err := re.s.Request(ctx, r.FileRecordID().FileBlockID())
	if err != nil {
		return err
	}

	defer b.Release()

	p := re.getPage(r, b)
	if re.skipLSN(r, p) {
		return nil
	}

	p.RecoverUpdate(r.FileRecordID().SlotID(), r.Data())
	return nil
}

func (re *Replay) recordClear(ctx context.Context, r *record.Record) error {
	b, err := re.s.Request(ctx, r.FileRecordID().FileBlockID())
	if err != nil {
		return err
	}

	defer b.Release()

	p := re.getPage(r, b)
	if re.skipLSN(r, p) {
		return nil
	}

	p.RecoverClear()
	return nil
}

func (re *Replay) recoverRollback() error {
	return nil
}

func (re *Replay) getPage(r *record.Record, b *storage.Buffer) *page.Page {
	if tablename.IsIndex(r.FileRecordID().FileBlockID().FileName().String()) {
		return page.NewFromSlice(b.Raw(), 0, span.FragmentHeaderSize)
	}

	return page.NewFromSlice(b.Raw(), node.HeaderSize, 0)
}

func (re *Replay) skipLSN(r *record.Record, p *page.Page) bool {
	return p.TrackingID() >= r.LogSequenceNumber().Value()
}
