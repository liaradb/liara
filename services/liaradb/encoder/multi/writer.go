package multi

import "io"

type Writer[W io.Writer] struct {
	writers []W
	current int
}

func NewWriter[W io.Writer](writers ...W) *Writer[W] {
	return &Writer[W]{
		writers: writers,
	}
}

func (wr *Writer[W]) Append(r W) {
	wr.writers = append(wr.writers, r)
}

func (wr *Writer[W]) Write(p []byte) (n int, err error) {
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
