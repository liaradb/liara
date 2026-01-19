package mempage

import "io"

type zeroHeader struct{}

var _ Serializer = zeroHeader{}

func (z zeroHeader) Read(io.Reader) error {
	return nil
}

func (z zeroHeader) Size() int {
	return 0
}

func (z zeroHeader) Write(io.Writer) error {
	return nil
}
