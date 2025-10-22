package value

import (
	"io"

	"github.com/liaradb/liaradb/raw"
)

type AggregateName string

func (an AggregateName) String() string         { return string(an) }
func (i AggregateName) Size() int               { return raw.StringSize(i) }
func (i AggregateName) Write(w io.Writer) error { return raw.WriteString(w, i) }
func (i *AggregateName) Read(r io.Reader) error { return raw.ReadString(r, i) }
