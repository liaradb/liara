package record

type Collection uint16

const CollectionSize = 2

const (
	CollectionSystem  Collection = 1
	CollectionRequest Collection = 2
	CollectionOutbox  Collection = 3
	CollectionEvent   Collection = 4
	CollectionValue   Collection = 5
)

func (c Collection) Size() int { return CollectionSize }

func (c Collection) String() string {
	switch c {
	case CollectionSystem:
		return "system"
	case CollectionRequest:
		return "request"
	case CollectionOutbox:
		return "outbox"
	case CollectionEvent:
		return "event"
	case CollectionValue:
		return "value"
	default:
		return ""
	}
}
