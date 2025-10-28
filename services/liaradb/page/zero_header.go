package page

import "io"

type ZeroHeader struct{}

var _ Serializer = ZeroHeader{}

func (z ZeroHeader) Read(io.Reader) error {
	return nil
}

func (z ZeroHeader) Size() int {
	return 0
}

func (z ZeroHeader) Write(io.Writer) error {
	return nil
}
