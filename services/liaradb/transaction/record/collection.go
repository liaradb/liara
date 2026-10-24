package record

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type Collection uint16

const CollectionSize = 2

const (
	CollectionSystem  Collection = 1
	CollectionTenant  Collection = 2
	CollectionRequest Collection = 3
	CollectionOutbox  Collection = 4
	CollectionEvent   Collection = 5
	CollectionValue   Collection = 6
	CollectionGraph   Collection = 7
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
	case CollectionTenant:
		return "tenant"
	case CollectionRequest:
		return "request"
	case CollectionOutbox:
		return "outbox"
	case CollectionEvent:
		return "event"
	case CollectionValue:
		return "value"
	case CollectionGraph:
		return "graph"
	default:
		return ""
	}
}
