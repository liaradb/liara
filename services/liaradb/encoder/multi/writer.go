package multi

import "io"

type Writer struct {
	writers []io.Writer
	current int
}

func NewWriter(writers ...io.Writer) *Writer {
	return &Writer{
		writers: writers,
	}
}

func (wr *Writer) Append(r io.Writer) {
	wr.writers = append(wr.writers, r)
}

func (wr *Writer) Write(p []byte) (n int, err error) {
	l := len(p)
	for wr.current < len(wr.writers) {
		w := wr.writers[wr.current]
		wn, err := w.Write(p[n:])
		n += wn
		if err != nil && err != io.ErrShortWrite {
			return n, err
		}
		if n >= l {
			return n, nil
		}
		wr.current++
	}
	return n, nil
}
