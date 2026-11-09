package span

import (
	"io"
	"slices"

	"github.com/liaradb/liaradb/encoder/multi"
	"github.com/liaradb/liaradb/encoder/page"
)

type Span struct {
	fragments []*Fragment
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

func (s *Span) Append(header []byte, data []byte) *Fragment {
	f := newFragment(header, data)
	s.fragments = append(s.fragments, f)
	return f
}

func (s *Span) Reverse() {
	slices.Reverse(s.fragments)
}

func (s *Span) InitIndexes() {
	c := len(s.fragments)
	l := c - 1
	for i, f := range s.fragments {
		f.setCount(int16(c))
		f.setIndex(int16(l - i))
	}
}

func (s Span) Read(p []byte) (n int, err error) {
	if !s.valid() {
		return 0, page.ErrInvalidCRC
	}

	reader := multi.NewReader(s.fragments...)
	return reader.Read(p)
}

func (s Span) Write(p []byte) (n int, err error) {
	writer := multi.NewWriter(s.fragments...)
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
