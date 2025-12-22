package leafnode

import (
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/storage"
)

type leafEntry struct {
	key      value.Key
	recordID storage.RecordLocator
}

func newLeafEntry(key value.Key, recordID storage.RecordLocator) leafEntry {
	return leafEntry{
		key:      key,
		recordID: recordID,
	}
}

func (le leafEntry) Key() value.Key                  { return le.key }
func (le leafEntry) RecordID() storage.RecordLocator { return le.recordID }
func (le leafEntry) Size() int                       { return le.key.Size() + storage.RecordIDSize }

func (le leafEntry) Write(data []byte) {
	data0 := le.recordID.Write(data)
	copy(data0, []byte(le.key))
}

func (le *leafEntry) Read(data []byte) {
	data0 := le.recordID.Read(data)
	le.key = value.Key(data0)
}
