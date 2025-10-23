package value

type EventName struct {
	baseString
}

func NewEventName(value string) EventName {
	return EventName{baseString(value)}
}
