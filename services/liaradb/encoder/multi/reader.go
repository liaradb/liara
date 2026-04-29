package multi

import "io"

type Reader struct {
	readers []io.Reader
}

func NewReader(readers ...io.Reader) *Reader {
	return &Reader{
		readers: readers,
	}
}

func (rd *Reader) Append(r io.Reader) {
	rd.readers = append(rd.readers, r)
}

func (rd *Reader) Read(p []byte) (n int, err error) {
	wn := 0
	for _, w := range rd.readers {
		wn, err = w.Read(p[n:])
		n += wn
		if err != nil && err != io.EOF {
			return n, err
		}
	}
	return n, err
}
