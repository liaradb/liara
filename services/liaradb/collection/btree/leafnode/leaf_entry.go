package leafnode

import (
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/storage/link"
)

type leafEntry struct {
	key      value.Key
	recordID link.RecordLocator
}

func newLeafEntry(key value.Key, recordID link.RecordLocator) leafEntry {
	return leafEntry{
		key:      key,
		recordID: recordID,
	}
}

func (le leafEntry) Key() value.Key               { return le.key }
func (le leafEntry) RecordID() link.RecordLocator { return le.recordID }
func (le leafEntry) Size() int                    { return le.key.Size() + link.RecordIDSize }

func (le leafEntry) Write(data []byte) {
	data0 := le.recordID.Write(data)
	le.key.Write(data0)
}

func (le *leafEntry) Read(data []byte) {
	data0 := le.recordID.Read(data)
	le.key.Read(data0)
}
