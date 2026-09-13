package span

import (
	"context"
	"io"
	"slices"

	"github.com/liaradb/liaradb/encoder/multi"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage/link"
)

type Span struct {
	l         Log
	fragments []*Fragment
}

type Log interface {
	Append(context.Context, link.SlotID, []byte) (logpage.LogSequenceNumber, error)
}

func New(l Log) *Span {
	return &Span{
		l: l,
	}
}

func (s Span) Length() (l int) {
	for _, f := range s.fragments {
		l += f.length()
	}
	return
}

func (s Span) valid() bool {
	for _, f := range s.fragments {
		if !f.valid() {
			return false
		}
	}
	return true
}

// TODO: Ensure fragments are sorted by BlockID
func (s *Span) Append(b BufferPage, sid link.SlotID, header []byte, data []byte) *Fragment {
	f := newFragment(s.l, b, sid, header, data)
	s.fragments = append(s.fragments, f)
	return f
}

func (s *Span) Reverse() {
	slices.Reverse(s.fragments)
}

func (s *Span) InitIndexes() {
	if len(s.fragments) < 2 {
		return
	}

	for i, f := range s.fragments[:len(s.fragments)-1] {
		next := s.fragments[i+1]
		f.setNextPosition(next.p.BlockID().Position())
		// f.setNextSlotID()
	}
}

// TODO: Can we do this without creating a new reader?
func (s Span) Read(p []byte) (n int, err error) {
	if !s.valid() {
		return 0, page.ErrInvalidCRC
	}

	readers := make([]io.Reader, 0, len(s.fragments))
	for _, f := range s.fragments {
		readers = append(readers, f)
	}

	reader := multi.NewReader(readers...)
	return reader.Read(p)
}

// TODO: Can we do this without creating a new writer?
func (s Span) Write(p []byte) (n int, err error) {
	writers := make([]io.Writer, 0, len(s.fragments))
	for _, f := range s.fragments {
		writers = append(writers, f)
	}

	writer := multi.NewWriter(writers...)
	return writer.Write(p)
}

func (s Span) SeekStart() error {
	for _, s := range s.fragments {
		if _, err := s.buffer.Seek(0, io.SeekStart); err != nil {
			return err
		}
	}

	return nil
}

func (s Span) Commit() {
	for _, f := range s.fragments {
		f.commit()
	}
}
