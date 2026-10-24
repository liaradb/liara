package replay

import (
	"context"

	"github.com/liaradb/liaradb/collection"
	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/transaction/record"
)

type Replay struct {
	s           *storage.Storage
	collections *collection.Collections
	log         *log.Log
}

func NewReplay(
	s *storage.Storage,
	collections *collection.Collections,
	log *log.Log,
) *Replay {
	return &Replay{
		s:           s,
		collections: collections,
		log:         log,
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
		return re.recoverCheckpoint(ctx, r)
	case record.ActionStart:
		return re.recoverStart(ctx, r)
	case record.ActionCommit:
		return re.recoverCommit(ctx, r)
	case record.ActionInsert:
		return re.recoverInsert(ctx, r)
	case record.ActionRemove:
		return re.recoverRemove(ctx, r)
	case record.ActionRollback:
		return re.recoverRollback(ctx, r)
	case record.ActionUpdate:
		return re.recoverUpdate(ctx, r)
	default:
		return ErrActionUnknown
	}
}

func (re *Replay) recoverCheckpoint(ctx context.Context, r *record.Record) error {
	return nil
}

func (re *Replay) recoverStart(ctx context.Context, r *record.Record) error {
	return nil
}

func (re *Replay) recoverCommit(ctx context.Context, r *record.Record) error {
	return nil
}

func (re *Replay) recoverInsert(ctx context.Context, r *record.Record) error {
	b, err := re.s.Request(ctx, r.RecordID().BlockID())
	if err != nil {
		return err
	}

	defer b.Release()

	p := re.getPage(r, b)
	p.RecoverInsert(r.RecordID().SlotID(), r.Data())
	return nil
}

func (re *Replay) recoverRemove(ctx context.Context, r *record.Record) error {
	b, err := re.s.Request(ctx, r.RecordID().BlockID())
	if err != nil {
		return err
	}

	defer b.Release()

	return nil
}

func (re *Replay) recoverUpdate(ctx context.Context, r *record.Record) error {
	b, err := re.s.Request(ctx, r.RecordID().BlockID())
	if err != nil {
		return err
	}

	defer b.Release()

	p := re.getPage(r, b)
	p.RecoverUpdate(r.RecordID().SlotID(), r.Data())
	return nil
}

func (re *Replay) recoverRollback(ctx context.Context, r *record.Record) error {
	return nil
}

func (re *Replay) getPage(r *record.Record, b *storage.Buffer) *page.Page {
	if tablename.IsIndex(r.RecordID().BlockID().FileName().String()) {
		return page.NewFromSlice(b.Raw(), 0, span.FragmentHeaderSize)
	}

	return page.NewFromSlice(b.Raw(), node.HeaderSize, 0)
}
