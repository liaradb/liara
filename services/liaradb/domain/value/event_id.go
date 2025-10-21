package value

import (
	"io"

	"github.com/google/uuid"
)

type EventID string

func (i EventID) String() string { return string(i) }

func NewEventID() EventID {
	return EventID(uuid.NewString())
}

const EventIDSize = 36

func (i EventID) Size() int { return EventIDSize }

func (i EventID) Write(w io.Writer) error {
	if n, err := w.Write([]byte(i)); err != nil {
		return err
	} else if n < len(i) {
		return io.ErrShortWrite
	}

	return nil
}

func (i *EventID) Read(r io.Reader) error {
	d := make([]byte, EventIDSize)
	if n, err := r.Read(d); err != nil {
		return err
	} else if n < EventIDSize {
		return io.EOF
	}

	*i = EventID(d)

	return nil
}
