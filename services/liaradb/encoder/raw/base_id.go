package raw

import (
	"io"

	"github.com/google/uuid"
)

type baseUUID = uuid.UUID

type BaseID struct {
	baseUUID
}

func NewBaseID() BaseID {
	return BaseID{uuid.New()}
}

func NewBaseIDFromString(value string) (BaseID, error) {
	if id, err := uuid.Parse(value); err != nil {
		return BaseID{}, err
	} else {
		return BaseID{id}, nil
	}
}

func (i BaseID) String() string { return i.baseUUID.String() }

const BaseIDSize = 16

func (i BaseID) Size() int { return len(i.baseUUID) }

func (i BaseID) Write(w io.Writer) error {
	if n, err := w.Write(i.baseUUID[:]); err != nil {
		return err
	} else if n < len(i.baseUUID) {
		return io.ErrShortWrite
	}

	return nil
}

func (i *BaseID) Read(r io.Reader) error {
	if n, err := r.Read(i.baseUUID[:]); err != nil {
		return err
	} else if n < BaseIDSize {
		return io.EOF
	}

	return nil
}
