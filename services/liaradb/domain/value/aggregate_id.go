package value

import (
	"io"

	"github.com/google/uuid"
	"github.com/liaradb/liaradb/raw"
)

type AggregateID string

func NewAggregateID() AggregateID {
	return AggregateID(uuid.NewString())
}

func (i AggregateID) String() string          { return string(i) }
func (i AggregateID) Size() int               { return raw.StringSize(i) }
func (i AggregateID) Write(w io.Writer) error { return raw.WriteString(w, i) }
func (i *AggregateID) Read(r io.Reader) error { return raw.ReadString(r, i) }
