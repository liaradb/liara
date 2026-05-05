package multi

import "io"

type Reader struct {
	readers []io.Reader
	current int
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
	l := len(p)
	for rd.current < len(rd.readers) {
		r := rd.readers[rd.current]
		wn, err := r.Read(p[n:])
		n += wn
		if err != nil && err != io.EOF {
			return n, err
		}
		if n >= l {
			return n, nil
		}
		rd.current++
	}
	return n, err
}
