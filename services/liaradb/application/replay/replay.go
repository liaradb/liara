package replay

import (
	"context"
	"fmt"

	"github.com/liaradb/liaradb/collection"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/record"
)

type Replay struct {
	collections *collection.Collections
	log         *recovery.Log
}

func NewReplay(
	collections *collection.Collections,
	log *recovery.Log,
) *Replay {
	return &Replay{
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
	fmt.Printf("recover: %v\n", r.Action())
	return nil
}

func (re *Replay) recoverCommit(ctx context.Context, r *record.Record) error {
	fmt.Printf("recover: %v\n", r.Action())
	return nil
}

func (re *Replay) recoverInsert(ctx context.Context, r *record.Record) error {
	switch r.Collection() {
	case record.CollectionEvent:
		return re.recoverEvent(ctx, r)
	default:
		fmt.Printf("recover: %v\n", r.Action())
		return nil
	}
}

func (re *Replay) recoverEvent(ctx context.Context, r *record.Record) error {
	var e entity.Event
	if err := e.Read(buffer.NewFromSlice(r.Data())); err != nil {
		return err
	}

	fmt.Printf("recover: %v: %v\n", r.Action(), e.AggregateID.String())
	tn := tablename.New(r.TenantID())
	return re.collections.EventLog.Append(ctx, tn, e.PartitionID, &e)
}

func (re *Replay) recoverRemove(ctx context.Context, r *record.Record) error {
	fmt.Printf("recover: %v\n", r.Action())
	return nil
}

func (re *Replay) recoverRollback(ctx context.Context, r *record.Record) error {
	fmt.Printf("recover: %v\n", r.Action())
	return nil
}

func (re *Replay) recoverUpdate(ctx context.Context, r *record.Record) error {
	switch r.Collection() {
	default:
		fmt.Printf("recover: %v\n", r.Action())
		return nil
	}
}
