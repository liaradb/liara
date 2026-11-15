package heap

import (
	"context"
	"io"

	"github.com/liaradb/liaradb/collection/heap/bufferpage"
	"github.com/liaradb/liaradb/collection/heap/record"
	"github.com/liaradb/liaradb/collection/heap/tip"
	"github.com/liaradb/liaradb/collection/log"
	"github.com/liaradb/liaradb/encoder/serializer"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type Heap struct {
	s *storage.Storage
}

func New(s *storage.Storage) *Heap {
	return &Heap{s: s}
}

func (h *Heap) Insert(
	ctx context.Context,
	l log.Log,
	fn link.FileName,
	w serializer.SizedSerializer,
) (link.RecordID, error) {
	t := tip.NewTip(h.s, fn)

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

	return s.RecordID(), nil
}

func (h *Heap) Get(
	ctx context.Context,
	fn link.FileName,
	rid link.RecordID,
) (*record.Span, error) {
	fbid := fn.FileBlockID(rid.Block())
	b, err := h.s.Request(ctx, fbid)
	if err != nil {
		return nil, err
	}

	s := record.New()
	p := bufferpage.New(b)
	f := s.AppendSlot(p, rid.SlotID())

	for f.NextBlockID() != 0 {
		fbid.SetBlockID(f.NextBlockID())
		b, err := h.s.Request(ctx, fbid)
		if err != nil {
			s.Release()
			return nil, err
		}

		p := bufferpage.New(b)
		f = s.AppendSlot(p, link.SlotID(0))
	}

	return s, nil
}

func (h *Heap) Replace(
	ctx context.Context,
	l log.Log,
	fn link.FileName,
	rid link.RecordID,
	w serializer.SizedSerializer,
) error {
	s, err := h.Get(ctx, fn, rid)
	if err != nil {
		return err
	}

	defer s.Release()

	if size, length := w.Size(), s.Length(); size < length {
		// Extend Span
		return io.ErrShortWrite
	} else if size > length {
		// Return Fragments to Free Space Manager
	}

	if err := w.Write(s.Writer()); err != nil {
		return err
	}

	return s.CommitUpdate(ctx, l)
}

func (h *Heap) Delete(
	ctx context.Context,
	l log.Log,
	fn link.FileName,
	rid link.RecordID,
) error {
	s, err := h.Get(ctx, fn, rid)
	if err != nil {
		return err
	}

	defer s.Release()

	s.Delete()
	// Return fragments to Free Space Manager

	return s.CommitUpdate(ctx, l)
}
