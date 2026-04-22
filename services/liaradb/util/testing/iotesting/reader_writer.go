package iotesting

import (
	"bufio"
	"bytes"
)

func NewReaderWriter() (*bufio.Reader, *bufio.Writer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer),
		bufio.NewWriter(buffer)
}

func NewReaderBuffer() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
