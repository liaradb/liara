package raw

import (
	"io"
)

type BaseString string

func (s BaseString) String() string          { return string(s) }
func (s BaseString) Bytes() []byte           { return []byte(s) }
func (s BaseString) Length() int             { return len(s) }
func (s BaseString) Size() int               { return StringSize(s) }
func (s BaseString) Write(w io.Writer) error { return WriteString(w, s) }
func (s *BaseString) Read(r io.Reader) error { return ReadString(r, s) }

// TODO: Test this
func (s BaseString) WriteData(data []byte, colSize int) []byte {
	data[0] = byte(len(s))
	copy(data[1:colSize], []byte(s))
	clear(data[len(s):colSize])
	return data[colSize:]
}

// TODO: Test this
func (s *BaseString) ReadData(data []byte, colSize int) []byte {
	l := data[0]
	*s = BaseString(data[1 : 1+l])
	return data[1+l:]
}
