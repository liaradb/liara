package leafnode

import (
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/storage/link"
)

type leafEntry struct {
	key      key.Key
	recordID link.RecordLocator
}

func newLeafEntry(key key.Key, recordID link.RecordLocator) leafEntry {
	return leafEntry{
		key:      key,
		recordID: recordID,
	}
}

func (le leafEntry) Key() key.Key                 { return le.key }
func (le leafEntry) RecordID() link.RecordLocator { return le.recordID }
func (le leafEntry) Size() int                    { return le.key.Size() + link.RecordLocatorSize }

func (le leafEntry) Write(data []byte) bool {
	data0, ok := le.recordID.Write(data)
	if !ok {
		return false
	}

	return le.key.Write(data0)
}

func (le *leafEntry) Read(data []byte) bool {
	data0, ok := le.recordID.Read(data)
	if !ok {
		return false
	}

	return le.key.Read(data0)
}
