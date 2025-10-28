package raw

import (
	"io"
)

type BaseString string

func (i BaseString) String() string          { return string(i) }
func (i BaseString) Length() int             { return len(i) }
func (i BaseString) Size() int               { return StringSize(i) }
func (i BaseString) Write(w io.Writer) error { return WriteString(w, i) }
func (i *BaseString) Read(r io.Reader) error { return ReadString(r, i) }
