package span

import (
	"io"
	"slices"

	"github.com/liaradb/liaradb/encoder/multi"
	"github.com/liaradb/liaradb/recovery/record"
)

type Span struct {
	fragments []Fragment
}

func NewSpan(fragments ...Fragment) Span {
	return Span{
		fragments: fragments,
	}
}

func (s Span) Length() (l int64) {
	for _, f := range s.fragments {
		l += f.Length()
	}
	return
}

func (s *Span) Append(f Fragment) {
	s.fragments = append(s.fragments, f)
}

func (s *Span) Reverse() {
	slices.Reverse(s.fragments)
}

func (s *Span) InitIndexes() {
	c := len(s.fragments)
	l := c - 1
	for i, f := range s.fragments {
		f.SetCount(int16(c))
		f.SetIndex(int16(l - i))
	}
}

// TODO: Can we do this without creating a new reader?
func (s Span) Read(p []byte) (n int, err error) {
	readers := make([]io.Reader, 0, len(s.fragments))
	for _, f := range s.fragments {
		readers = append(readers, &f)
	}

	reader := multi.NewReader(readers...)
	return reader.Read(p)
}

// TODO: Can we do this without creating a new writer?
func (s Span) Write(p []byte) (n int, err error) {
	writers := make([]io.Writer, 0, len(s.fragments))
	for _, f := range s.fragments {
		writers = append(writers, &f)
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

func (s Span) ToRecord() (*record.Record, error) {
	rc := record.Record{}
	if err := rc.Read(s); err != nil {
		return nil, err
	}

	return &rc, nil
}
