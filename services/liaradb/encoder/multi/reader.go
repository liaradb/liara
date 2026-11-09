package multi

import "io"

type Reader[R io.Reader] struct {
	readers []R
	current int
}

func NewReader[R io.Reader](readers ...R) *Reader[R] {
	return &Reader[R]{
		readers: readers,
	}
}

func (rd *Reader[R]) Append(r R) {
	rd.readers = append(rd.readers, r)
}

func (rd *Reader[R]) Read(p []byte) (n int, err error) {
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
