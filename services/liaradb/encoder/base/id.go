package base

import (
	"io"

	"github.com/google/uuid"
)

type baseUUID = uuid.UUID

type ID struct {
	baseUUID
}

func NewID() ID {
	return ID{uuid.New()}
}

func NewIDFromString(value string) (ID, error) {
	if id, err := uuid.Parse(value); err != nil {
		return ID{}, err
	} else {
		return ID{id}, nil
	}
}

func (i ID) String() string { return i.baseUUID.String() }
func (i ID) Bytes() []byte  { return i.baseUUID[:] }

const BaseIDSize = 16

func (i ID) Size() int { return len(i.baseUUID) }

func (i ID) Write(w io.Writer) error {
	if n, err := w.Write(i.baseUUID[:]); err != nil {
		return err
	} else if n < len(i.baseUUID) {
		return io.ErrShortWrite
	}

	return nil
}

func (i *ID) Read(r io.Reader) error {
	if n, err := r.Read(i.baseUUID[:]); err != nil {
		return err
	} else if n < BaseIDSize {
		return io.EOF
	}

	return nil
}

func (i ID) WriteData(data []byte) []byte {
	copy(data[:16], i.baseUUID[:])
	return data[16:]
}

func (i *ID) ReadData(data []byte) []byte {
	copy(i.baseUUID[:], data[:16])
	return data[16:]
}
