package value

import (
	"io"

	"github.com/google/uuid"
)

type EventID [EventIDSize]byte

func (i EventID) String() string { return uuid.UUID(i).String() }

func NewEventID() EventID {
	return EventID(uuid.New())
}

const EventIDSize = 16

func (i EventID) Size() int { return EventIDSize }

func (i EventID) Write(w io.Writer) error {
	if n, err := w.Write(i[:]); err != nil {
		return err
	} else if n < len(i) {
		return io.ErrShortWrite
	}

	return nil
}

func (i *EventID) Read(r io.Reader) error {
	if n, err := r.Read(i[:]); err != nil {
		return err
	} else if n < EventIDSize {
		return io.EOF
	}

	return nil
}
