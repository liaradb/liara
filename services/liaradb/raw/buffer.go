package raw

import "io"

type Buffer struct {
	data   []byte
	cursor int64
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		data: make([]byte, size),
	}
}

type w interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt
}

var _ w = (*Buffer)(nil)

func (b *Buffer) Read(p []byte) (n int, err error) {
	if n = copy(p, b.data[b.cursor:]); n < len(p) {
		err = io.EOF
	}
	b.cursor += int64(n)
	return
}

func (b *Buffer) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, ErrUnderflow
	}

	b.cursor = off
	return b.Read(p)
}

func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	c := b.cursor
	switch whence {
	case io.SeekStart:
		c = offset
	case io.SeekCurrent:
		c += offset
	case io.SeekEnd:
		c = int64(len(b.data)) + offset
	default:
		c = offset
	}

	if c < 0 {
		return b.cursor, ErrUnderflow
	}

	b.cursor = c
	return b.cursor, nil
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	if n = copy(b.data[b.cursor:], p); n < len(p) {
		err = io.ErrShortWrite
	}
	b.cursor += int64(n)
	return
}

func (b *Buffer) WriteAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, ErrUnderflow
	}

	b.cursor = off
	return b.Write(p)
}
