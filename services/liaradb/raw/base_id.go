package raw

import (
	"io"

	"github.com/google/uuid"
)

type BaseID [BaseIDSize]byte

func (i BaseID) String() string { return uuid.UUID(i).String() }

func NewBaseID() BaseID {
	return BaseID(uuid.New())
}

const BaseIDSize = 16

func (i BaseID) Size() int { return BaseIDSize }

func (i BaseID) Write(w io.Writer) error {
	if n, err := w.Write(i[:]); err != nil {
		return err
	} else if n < len(i) {
		return io.ErrShortWrite
	}

	return nil
}

func (i *BaseID) Read(r io.Reader) error {
	if n, err := r.Read(i[:]); err != nil {
		return err
	} else if n < BaseIDSize {
		return io.EOF
	}

	return nil
}
