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

	r, err := t.Record(ctx, w.Size())
	if err != nil {
		return link.RecordID{}, err
	}

	defer r.Release()

	if err := w.Write(r.Writer()); err != nil {
		return link.RecordID{}, err
	}

	if err := r.CommitAppend(ctx, l); err != nil {
		return link.RecordID{}, err
	}

	return r.RecordID(), nil
}

func (h *Heap) Get(
	ctx context.Context,
	fn link.FileName,
	rid link.RecordID,
) (*record.Record, error) {
	fbid := fn.FileBlockID(rid.Block())
	b, err := h.s.Request(ctx, fbid)
	if err != nil {
		return nil, err
	}

	r := record.New()
	ok := false
	defer func() {
		if !ok {
			r.Release()
		}
	}()

	p := bufferpage.New(b)
	f := r.AppendSlot(p, rid.SlotID())

	for f.NextBlockID() != 0 {
		fbid.SetBlockID(f.NextBlockID())
		b, err := h.s.Request(ctx, fbid)
		if err != nil {
			return nil, err
		}

		p := bufferpage.New(b)
		f = r.AppendSlot(p, f.NextSlotID())
	}

	ok = true
	return r, nil
}

// TODO: Change to take a Record instead of Serializer to avoid second Get
func (h *Heap) Replace(
	ctx context.Context,
	l log.Log,
	fn link.FileName,
	rid link.RecordID,
	w serializer.SizedSerializer,
) error {
	r, err := h.Get(ctx, fn, rid)
	if err != nil {
		return err
	}

	defer r.Release()

	if size, length := w.Size(), r.Length(); size < length {
		// Extend Record
		return io.ErrShortWrite
	} else if size > length {
		// Return Fragments to Free Space Manager
	}

	if err := w.Write(r.Writer()); err != nil {
		return err
	}

	return r.CommitUpdate(ctx, l)
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
