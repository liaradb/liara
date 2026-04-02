package cursor

import "io"

type multiWriter struct {
	writers []io.Writer
}

func newMultiWriter(writers ...io.Writer) *multiWriter {
	return &multiWriter{
		writers: writers,
	}
}

func (mw *multiWriter) Write(p []byte) (int, error) {
	n := 0
	for _, w := range mw.writers {
		wn, err := w.Write(p[n:])
		n += wn
		if err != nil && err != io.ErrShortWrite {
			return n, err
		}
	}
	return n, nil
}
