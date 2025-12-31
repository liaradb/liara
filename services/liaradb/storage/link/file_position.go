package link

import (
	"fmt"
	"io"

	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/encoder/raw"
)

type FilePosition int64

const FilePositionSize = 8

func (p FilePosition) Value() int64   { return int64(p) }
func (FilePosition) Size() int        { return FilePositionSize }
func (p FilePosition) String() string { return fmt.Sprintf("%v", p.Value()) }

func (p FilePosition) Offset(bufferSize int64) page.Offset {
	return page.Offset(p) * page.Offset(bufferSize)
}

func (p FilePosition) Write(w io.Writer) error {
	return raw.WriteInt64(w, p)
}

func (p *FilePosition) Read(r io.Reader) error {
	return raw.ReadInt64(r, p)
}
