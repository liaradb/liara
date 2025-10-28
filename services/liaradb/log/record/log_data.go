package record

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type LogData struct {
	data []byte
}

const LogDataHeaderSize = 4

func NewLogData(data []byte) *LogData {
	return &LogData{
		data: data,
	}
}

func (ld *LogData) Bytes() []byte           { return ld.data }
func (ld *LogData) Length() int             { return len(ld.data) }
func (ld *LogData) Size() int               { return raw.HeaderSize + len(ld.data) }
func (ld *LogData) Write(w io.Writer) error { return raw.Write(w, ld.data) }
func (ld *LogData) Read(r io.Reader) error  { return raw.Read(r, &ld.data) }
