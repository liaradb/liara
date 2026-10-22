package span

import (
	"io"
	"slices"

	"github.com/liaradb/liaradb/collection/bufferpage"
	"github.com/liaradb/liaradb/encoder/multi"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type Span struct {
	fragments []*Fragment
	buffers   []BufferPage
}

func New() *Span {
	return &Span{}
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

func (s *Span) AppendSlot(b *storage.Buffer, sid link.SlotID) *Fragment {
	p := bufferpage.New(b, FragmentHeaderSize)
	h, d := p.Slot(sid)
	return s.Append(p, sid, h, d)
}

func (s *Span) Append(b BufferPage, sid link.SlotID, header []byte, data []byte) *Fragment {
	s.buffers = append(s.buffers, b)
	f := newFragment(b, sid, header, data)
	s.fragments = append(s.fragments, f)
	return f
}

func (s *Span) InitIndexes() {
	if len(s.fragments) < 2 {
		return
	}

	s.sortFragments()

	for i, f := range s.fragments[:len(s.fragments)-1] {
		next := s.fragments[i+1]
		f.setNextPosition(next.p.BlockID().Position())
		f.setNextSlotID(0) // TODO: Is this necessary?
	}
}

// Sort to ensure writes are sequential, and prevent deadlocks
func (s *Span) sortFragments() {
	slices.SortFunc(s.fragments, func(a, b *Fragment) int {
		return int(a.p.BlockID().Position() - b.p.BlockID().Position())
	})
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
func (s *Span) Writer() io.Writer {
	writers := make([]io.Writer, 0, len(s.fragments))
	for _, f := range s.fragments {
		f.prepareWrite()
		writers = append(writers, f)
	}

	return multi.NewWriter(writers...)
}

// TODO: This is unused
func (s Span) SeekStart() {
	for _, f := range s.fragments {
		f.buffer.SeekStart()
	}
}

func (s Span) Commit() {
	for _, f := range s.fragments {
		f.commit()
	}
}

func (s Span) LogAppend(l Log) error {
	for _, f := range s.fragments {
		if err := f.logAppend(l); err != nil {
			return err
		}
	}

	return nil
}

func (s Span) LogUpdate(l Log) error {
	for _, f := range s.fragments {
		if err := f.logUpdate(l); err != nil {
			return err
		}
	}

	return nil
}

func (s *Span) Release() {
	for _, b := range s.buffers {
		b.Release()
	}
}

func (s *Span) Bytes() ([]byte, error) {
	// Read Span
	buffer := make([]byte, s.Length())
	if _, err := s.Read(buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

func (s *Span) BytesAndRelease() ([]byte, error) {
	defer s.Release()

	return s.Bytes()
}
