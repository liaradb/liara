package base

import (
	"io"
	"uuid"
)

const BaseIDSize = 16

type baseUUID = uuid.UUID

type ID struct {
	baseUUID
}

func NewID() ID {
	return ID{uuid.NewV7()}
}

func NewIDFromUUID(id uuid.UUID) ID {
	return ID{id}
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

func (i ID) Size() int { return BaseIDSize }

func (i ID) Write(w io.Writer) error {
	if _, err := w.Write(i.baseUUID[:]); err != nil {
		return err
	}

	return nil
}

func (i *ID) Read(r io.Reader) error {
	if _, err := r.Read(i.baseUUID[:]); err != nil {
		return err
	}

	return nil
}
