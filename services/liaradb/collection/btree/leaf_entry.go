package btree

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

// TODO: Test this
type LeafEntry struct {
	key      Key
	recordID RecordID
}

func newLeafEntry(key Key, recordID RecordID) LeafEntry {
	return LeafEntry{
		key:      key,
		recordID: recordID,
	}
}

func (le LeafEntry) Size() int { return le.key.Size() + RecordIDSize }

func (le LeafEntry) Write(w io.Writer) error {
	return raw.WriteAll(w,
		le.key,
		le.recordID,
	)
}

func (le *LeafEntry) Read(r io.Reader) error {
	return raw.ReadAll(r,
		&le.key,
		&le.recordID)
}
