package raw

import (
	"encoding/binary"
	"io"
)

const HeaderSize = 4

func StringSize[S ~string](s S) int {
	return HeaderSize + len(s)
}

func Write(w io.Writer, value []byte) error {
	if err := WriteInt32(w, uint32(len(value))); err != nil {
		return err
	}

	if _, err := w.Write(value); err != nil {
		return err
	}

	return nil
}

func Read(r io.Reader, value *[]byte) error {
	l, err := readLength(r)
	if err != nil {
		return err
	}

	v := *value
	if l > uint32(len(v)) {
		*value = make([]byte, l)
	} else {
		if v == nil {
			v = make([]byte, 0)
		} else {
			clear(v)
			v = v[:l]
		}
		*value = v
	}

	return readData(r, *value)
}

func WriteString[S ~string](w io.Writer, value S) error {
	if err := WriteInt32(w, uint32(len(value))); err != nil {
		return err
	}

	if _, err := w.Write([]byte(value)); err != nil {
		return err
	}

	return nil
}

func ReadString[S ~string](r io.Reader, s *S) error {
	l, err := readLength(r)
	if err != nil {
		return err
	}

	b := make([]byte, l)
	if err := readData(r, b); err != nil {
		return err
	}

	*s = S(b)
	return nil
}

func readLength(r io.Reader) (uint32, error) {
	var l uint32
	err := ReadInt32(r, &l)
	return l, err
}

func readData(r io.Reader, b []byte) error {
	_, err := r.Read(b)
	return err
}

func WriteInt8[T ~uint8 | ~int8](w io.Writer, v T) error {
	d := [1]byte{}
	d[0] = byte(v)
	_, err := w.Write(d[:])
	return err
}

func WriteInt16[T ~uint16 | ~int16](w io.Writer, v T) error {
	d := [2]byte{}
	binary.BigEndian.PutUint16(d[:], uint16(v))
	_, err := w.Write(d[:])
	return err
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

func ReadInt8[T ~uint8 | ~int8](r io.Reader, v *T) error {
	d := [1]byte{}
	if err := readToSlice(r, d[:]); err != nil {
		return err
	}

	*v = T(d[0])
	return nil
}

func ReadInt16[T ~uint16 | ~int16](r io.Reader, v *T) error {
	d := [2]byte{}
	if err := readToSlice(r, d[:]); err != nil {
		return err
	}

	*v = T(binary.BigEndian.Uint16(d[:]))
	return nil
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
	if _, err := r.Read(d); err != nil {
		return err
	}

	return nil
}
