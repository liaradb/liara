package eventlog

type Event struct {
	id EventID
}

func (e *Event) ID() EventID { return e.id }
