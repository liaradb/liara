package raw

import (
	"encoding/binary"
	"io"
)

const HeaderSize = 4

// TODO: Can we rename this Size?
func StringSize[S ~string](s S) int {
	return HeaderSize + len(s)
}

func Write(w io.Writer, value []byte) error {
	if err := binary.Write(w, binary.BigEndian, uint32(len(value))); err != nil {
		return err
	}

	if n, err := w.Write([]byte(value)); err != nil {
		return err
	} else if n < len(value) {
		return io.ErrShortWrite
	}

	return nil
}

// TODO: Can we reuse value?
func Read(r io.Reader, value *[]byte) error {
	l, err := readLength(r)
	if err != nil {
		return err
	}

	b, err := readData(r, l)
	if err != nil {
		return err
	}

	*value = b
	return nil
}

func WriteString[S ~string](w io.Writer, value S) error {
	if err := binary.Write(w, binary.BigEndian, uint32(len(value))); err != nil {
		return err
	}

	if n, err := w.Write([]byte(value)); err != nil {
		return err
	} else if n < len(value) {
		return io.ErrShortWrite
	}

	return nil
}

func ReadString[S ~string](r io.Reader, s *S) error {
	l, err := readLength(r)
	if err != nil {
		return err
	}

	b, err := readData(r, l)
	if err != nil {
		return err
	}

	*s = S(b)
	return nil
}

func readLength(r io.Reader) (uint32, error) {
	var l uint32
	err := binary.Read(r, binary.BigEndian, &l)
	return l, err
}

func readData(r io.Reader, l uint32) ([]byte, error) {
	b := make([]byte, l)

	if _, err := r.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

func WriteInt32[T ~uint32 | ~int32](w io.Writer, v T) error {
	d := [4]byte{}
	binary.BigEndian.PutUint32(d[:], uint32(v))
	_, err := w.Write(d[:])
	return err
}

func WriteInt64[T ~uint64 | ~int64](w io.Writer, v T) error {
	d := [8]byte{}
	binary.BigEndian.PutUint64(d[:], uint64(v))
	_, err := w.Write(d[:])
	return err
}

func ReadInt32[T ~uint32 | ~int32](r io.Reader, v *T) error {
	d := [4]byte{}
	if err := readToSlice(r, d[:]); err != nil {
		return err
	}

	*v = T(binary.BigEndian.Uint32(d[:]))
	return nil
}

func ReadInt64[T ~uint64 | ~int64](r io.Reader, v *T) error {
	d := [8]byte{}
	if err := readToSlice(r, d[:]); err != nil {
		return err
	}

	*v = T(binary.BigEndian.Uint64(d[:]))
	return nil
}

func readToSlice(r io.Reader, d []byte) error {
	if n, err := r.Read(d); err != nil {
		return err
	} else if n < len(d) {
		return io.ErrUnexpectedEOF
	}

	return nil
}
