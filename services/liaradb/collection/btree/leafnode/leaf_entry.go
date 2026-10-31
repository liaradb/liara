package leafnode

import (
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/storage/link"
)

type leafEntry struct {
	key      key.Key
	recordID link.RecordID
}

func newLeafEntry(key key.Key, recordID link.RecordID) leafEntry {
	return leafEntry{
		key:      key,
		recordID: recordID,
	}
}

func (le leafEntry) Key() key.Key            { return le.key }
func (le leafEntry) RecordID() link.RecordID { return le.recordID }
func (le leafEntry) Size() int               { return le.key.Size() + link.RecordIDSize }

func (le leafEntry) Write(data []byte) {
	data0 := le.recordID.WriteData(data)
	le.key.Write(data0)
}

func (le *leafEntry) Read(data []byte) {
	data0 := le.recordID.ReadData(data)
	le.key.Read(data0)
}
