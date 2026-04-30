package record

import (
	"io"

	"github.com/liaradb/liaradb/encoder/multi"
)

type Span struct {
	fragments []Fragment
	reader    *multi.Reader
	writer    *multi.Writer
}

func NewSpan(fragments ...Fragment) Span {
	readers := make([]io.Reader, 0, len(fragments))
	for _, f := range fragments {
		readers = append(readers, &f)
	}

	writers := make([]io.Writer, 0, len(fragments))
	for _, f := range fragments {
		writers = append(writers, &f)
	}

	return Span{
		fragments: fragments,
		reader:    multi.NewReader(readers...),
		writer:    multi.NewWriter(writers...),
	}
}

func (s Span) Read(p []byte) (n int, err error) {
	return s.reader.Read(p)
}

func (s Span) Write(p []byte) (n int, err error) {
	return s.writer.Write(p)
}

func (s Span) SeekStart() error {
	for _, s := range s.fragments {
		if _, err := s.buffer.Seek(0, io.SeekStart); err != nil {
			return err
		}
	}

	return nil
}
