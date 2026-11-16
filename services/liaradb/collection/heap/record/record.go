package record

import (
	"context"
	"io"
	"slices"

	"github.com/liaradb/liaradb/collection/log"
	"github.com/liaradb/liaradb/encoder/multi"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage/link"
)

type Record struct {
	fragments []*Fragment
	buffers   []BufferPage
}

func New() *Record {
	return &Record{}
}

func (r *Record) RecordID() link.RecordID {
	return r.fragments[0].RecordID()
}

func (r *Record) Length() (l int) {
	for _, f := range r.fragments {
		l += f.length()
	}
	return
}

func (r *Record) valid() bool {
	for _, f := range r.fragments {
		if !f.valid() {
			return false
		}
	}
	return true
}

func (r *Record) AppendSlot(b BufferPage, sid link.SlotID) *Fragment {
	// Append first in case of panic
	r.buffers = append(r.buffers, b)
	h, d := b.Slot(sid)
	return r.appendFragment(b, sid, h, d)
}

func (r *Record) Append(b BufferPage, sid link.SlotID, header []byte, data []byte) *Fragment {
	r.buffers = append(r.buffers, b)
	return r.appendFragment(b, sid, header, data)
}

func (r *Record) appendFragment(b BufferPage, sid link.SlotID, header []byte, data []byte) *Fragment {
	f := newFragment(b, sid, header, data)
	r.fragments = append(r.fragments, f)
	return f
}

func (r *Record) InitIndexes() {
	if len(r.fragments) < 2 {
		return
	}

	r.sortFragments()

	for i, f := range r.fragments[:len(r.fragments)-1] {
		next := r.fragments[i+1]
		f.setNextBlockID(next.p.FileBlockID().BlockID())
		f.setNextSlotID(next.slotID)
	}
}

// Sort to ensure writes are sequential, and prevent deadlocks
func (r *Record) sortFragments() {
	slices.SortFunc(r.fragments, func(a, b *Fragment) int {
		return int(a.p.FileBlockID().BlockID() - b.p.FileBlockID().BlockID())
	})
}

func (r Record) Read(p []byte) (n int, err error) {
	if !r.valid() {
		return 0, page.ErrInvalidCRC
	}

	reader := multi.NewReader(r.fragments...)
	return reader.Read(p)
}

func (r *Record) Writer() io.Writer {
	for _, f := range r.fragments {
		f.prepareWrite()
	}

	return multi.NewWriter(r.fragments...)
}

// Testing only
func (r Record) SeekStart() {
	for _, f := range r.fragments {
		f.buffer.SeekStart()
	}
}

func (r Record) CommitAppend(ctx context.Context, l log.Log) error {
	r.commit()
	return r.logAppend(ctx, l)
}

func (r Record) CommitUpdate(ctx context.Context, l log.Log) error {
	r.commit()
	return r.logUpdate(ctx, l)
}

func (r Record) commit() {
	for _, f := range r.fragments {
		f.commit()
	}
}

func (r Record) logAppend(ctx context.Context, l log.Log) error {
	for _, f := range r.fragments {
		if err := f.logAppend(ctx, l); err != nil {
			return err
		}
	}

	return nil
}

func (r Record) logUpdate(ctx context.Context, l log.Log) error {
	for _, f := range r.fragments {
		if err := f.logUpdate(ctx, l); err != nil {
			return err
		}
	}

	return nil
}

func (r *Record) Release() {
	for _, b := range r.buffers {
		b.Release()
	}
}

func (r *Record) Bytes() ([]byte, error) {
	// Read Record
	buffer := make([]byte, r.Length())
	if _, err := r.Read(buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

func (r *Record) BytesAndRelease() ([]byte, error) {
	defer r.Release()

	return r.Bytes()
}

func (r *Record) Delete() {
	for _, f := range r.fragments {
		f.Delete()
	}
}
