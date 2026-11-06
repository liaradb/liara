package segment

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

// TODO: Should this be moved to link?
type PageID uint64

const PageIDSize = 8

// TODO: This is private
func newActivePageIDFromSize(size int64, pageSize int16) PageID {
	if pageSize == 0 || size == 0 {
		return 0
	}
	ps := int64(pageSize)
	pid := size / ps
	if size%ps == 0 {
		pid--
	}
	return PageID(pid)
}

func (pid PageID) Value() uint64 {
	return uint64(pid)
}

func (pid PageID) Position(pageSize int16) int64 {
	return int64(pid) * int64(pageSize)
}

func (pid PageID) Write(w io.Writer) error {
	return raw.WriteInt64(w, pid)
}

func (pid *PageID) Read(r io.Reader) error {
	return raw.ReadInt64(r, pid)
}
