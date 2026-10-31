package graph

import "github.com/liaradb/liaradb/storage/link"

type EdgeProperty struct {
	id   link.RecordID
	next link.RecordID
}

const EdgePropertySize = link.RecordIDSize +
	link.RecordIDSize

func NewEdgeProperty(
	id link.RecordID,
	next link.RecordID,
) *EdgeProperty {
	return &EdgeProperty{
		id:   id,
		next: next,
	}
}

func (ep *EdgeProperty) ID() link.RecordID   { return ep.id }
func (ep *EdgeProperty) Next() link.RecordID { return ep.next }
func (ep *EdgeProperty) Size() int           { return EdgePropertySize }
