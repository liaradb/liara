package cursor

import "io"

type multiReader struct {
	readers []io.Reader
}

func newMultiReader(readers ...io.Reader) *multiReader {
	return &multiReader{
		readers: readers,
	}
}

func (mw *multiReader) Read(p []byte) (n int, err error) {
	wn := 0
	for _, w := range mw.readers {
		wn, err = w.Read(p[n:])
		n += wn
		if err != nil && err != io.EOF {
			return n, err
		}
	}
	return n, err
}
