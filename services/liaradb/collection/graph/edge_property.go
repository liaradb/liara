package graph

import "github.com/liaradb/liaradb/storage/link"

type EdgeProperty struct {
	id   link.RecordLocator
	next link.RecordLocator
}

const EdgePropertySize = link.RecordLocatorSize +
	link.RecordLocatorSize

func NewEdgeProperty(
	id link.RecordLocator,
	next link.RecordLocator,
) *EdgeProperty {
	return &EdgeProperty{
		id:   id,
		next: next,
	}
}

func (ep *EdgeProperty) ID() link.RecordLocator   { return ep.id }
func (ep *EdgeProperty) Next() link.RecordLocator { return ep.next }
func (ep *EdgeProperty) Size() int                { return EdgePropertySize }
