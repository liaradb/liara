package raw

import "io"

// TODO: Potentially use io.OffsetWriter
type Buffer struct {
	data   []byte
	cursor int64
}

func NewBuffer(size int64) *Buffer {
	return &Buffer{
		data: make([]byte, size),
	}
}

func NewBufferFromSlice(data []byte) *Buffer {
	return &Buffer{
		data: data,
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

func (b *Buffer) Clear() {
	clear(b.data)
	b.cursor = 0
}

// TODO: Test this
func (b *Buffer) ClearAfter(n int) {
	clear(b.data[n:])
}

func (b *Buffer) Reset(data []byte) {
	b.data = data
	b.cursor = 0
}

func (b *Buffer) Bytes() []byte { return b.data }
func (b *Buffer) Length() int64 { return int64(len(b.data)) }

func (b *Buffer) Read(p []byte) (n int, err error) {
	if b.cursor > int64(len(b.data)) {
		return 0, io.EOF
	}

	// TODO: Test this case
	if b.cursor == int64(len(b.data)) {
		return 0, nil
	}

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
	if b.cursor >= int64(len(b.data)) {
		return 0, io.ErrShortWrite
	}

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
