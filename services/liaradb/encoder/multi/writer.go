package multi

import "io"

type Writer struct {
	writers []io.Writer
}

func NewWriter(writers ...io.Writer) *Writer {
	return &Writer{
		writers: writers,
	}
}

func (w *Writer) Append(r io.Writer) {
	w.writers = append(w.writers, r)
}

func (w *Writer) Write(p []byte) (int, error) {
	n := 0
	for _, w := range w.writers {
		wn, err := w.Write(p[n:])
		n += wn
		if err != nil && err != io.ErrShortWrite {
			return n, err
		}
	}
	return n, nil
}
