package record

import (
	"context"

	"github.com/liaradb/liaradb/collection/record/tip"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/serializer"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/log"
)

type RecordCollection struct {
	s *storage.Storage
	l *log.Log
}

func NewRecordCollection(s *storage.Storage, l *log.Log) *RecordCollection {
	return &RecordCollection{
		s: s,
		l: l,
	}
}

func (rc *RecordCollection) Insert(
	ctx context.Context,
	tx span.Transaction,
	tn tablename.TableName,
	w serializer.SizedSerializer,
) (link.RecordID, error) {
	t := tip.NewTip(rc.s, tn.Collection())

	s, err := t.Span(ctx, w.Size())
	if err != nil {
		return link.RecordID{}, err
	}

	defer s.Release()

	if err := w.Write(s.Writer()); err != nil {
		return link.RecordID{}, err
	}

	if err := s.CommitAppend(ctx, tx); err != nil {
		return link.RecordID{}, err
	}

	return t.RecordID(), nil
}

func (rc *RecordCollection) Get(
	ctx context.Context,
	tn tablename.TableName,
	rid link.RecordID,
) (*span.Span, error) {
	fbid := tn.Collection().FileBlockID(rid.Block())
	b, err := rc.s.Request(ctx, fbid)
	if err != nil {
		return nil, err
	}

	s := span.New()
	f := s.AppendSlot(b, rid.SlotID())

	for f.NextBlockID() != 0 {
		fbid.SetBlockID(f.NextBlockID())
		b, err := rc.s.Request(ctx, fbid)
		if err != nil {
			s.Release()
			return nil, err
		}

		f = s.AppendSlot(b, link.SlotID(0))
	}

	return s, nil
}
