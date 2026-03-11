package serializer

import "io"

type Deserializer interface {
	Read(r io.Reader) error
}

func ReadAll(r io.Reader, ders ...Deserializer) error {
	for _, d := range ders {
		if err := d.Read(r); err != nil {
			return err
		}
	}

	return nil
}

type Serializer interface {
	Write(w io.Writer) error
}

func WriteAll(w io.Writer, sers ...Serializer) error {
	for _, s := range sers {
		if err := s.Write(w); err != nil {
			return err
		}
	}

	return nil
}
