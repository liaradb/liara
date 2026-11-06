package record

import (
	"context"

	"github.com/liaradb/liaradb/collection/bufferpage"
	"github.com/liaradb/liaradb/collection/record/tip"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/encoder/serializer"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type RecordCollection struct {
	s *storage.Storage
}

func New(s *storage.Storage) *RecordCollection {
	return &RecordCollection{s: s}
}

func (rc *RecordCollection) Insert(
	ctx context.Context,
	l span.Log,
	fn link.FileName,
	w serializer.SizedSerializer,
) (link.RecordID, error) {
	t := tip.NewTip(rc.s, fn)

	s, err := t.Span(ctx, w.Size())
	if err != nil {
		return link.RecordID{}, err
	}

	defer s.Release()

	if err := w.Write(s.Writer()); err != nil {
		return link.RecordID{}, err
	}

	if err := s.CommitAppend(ctx, l); err != nil {
		return link.RecordID{}, err
	}

	return t.RecordID(), nil
}

func (rc *RecordCollection) Get(
	ctx context.Context,
	fn link.FileName,
	rid link.RecordID,
) (*span.Span, error) {
	fbid := fn.FileBlockID(rid.Block())
	b, err := rc.s.Request(ctx, fbid)
	if err != nil {
		return nil, err
	}

	s := span.New()
	p := bufferpage.New(b)
	f := s.AppendSlot(p, rid.SlotID())

	for f.NextBlockID() != 0 {
		fbid.SetBlockID(f.NextBlockID())
		b, err := rc.s.Request(ctx, fbid)
		if err != nil {
			s.Release()
			return nil, err
		}

		p := bufferpage.New(b)
		f = s.AppendSlot(p, link.SlotID(0))
	}

	return s, nil
}

func (rc *RecordCollection) Replace(
	ctx context.Context,
	l span.Log,
	fn link.FileName,
	rid link.RecordID,
	w serializer.SizedSerializer,
) error {
	s, err := rc.Get(ctx, fn, rid)
	if err != nil {
		return err
	}

	defer s.Release()

	if w.Size() > s.Length() {
		// TODO: Verify data can fit
		panic("support resizing")
	}

	if err := w.Write(s.Writer()); err != nil {
		return err
	}

	return s.CommitUpdate(ctx, l)
}
