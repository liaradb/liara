package record

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

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

func (a Collection) Write(w io.Writer) error {
	return raw.WriteInt16(w, a)
}

func (a *Collection) Read(r io.Reader) error {
	return raw.ReadInt16(r, a)
}

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
