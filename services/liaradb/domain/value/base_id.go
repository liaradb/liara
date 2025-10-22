package value

import (
	"io"

	"github.com/google/uuid"
	"github.com/liaradb/liaradb/raw"
)

type baseID string

func newBaseID() baseID {
	return baseID(uuid.NewString())
}

func (i baseID) String() string          { return string(i) }
func (i baseID) Length() int             { return len(i) }
func (i baseID) Size() int               { return raw.StringSize(i) }
func (i baseID) Write(w io.Writer) error { return raw.WriteString(w, i) }
func (i *baseID) Read(r io.Reader) error { return raw.ReadString(r, i) }
