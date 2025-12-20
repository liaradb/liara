package leafnode

import "github.com/liaradb/liaradb/collection/btree/key"

type leafEntry struct {
	key      key.Key
	recordID RecordID
}

func newLeafEntry(key key.Key, recordID RecordID) leafEntry {
	return leafEntry{
		key:      key,
		recordID: recordID,
	}
}

func (le leafEntry) Key() key.Key       { return le.key }
func (le leafEntry) RecordID() RecordID { return le.recordID }
func (le leafEntry) Size() int          { return le.key.Size() + RecordIDSize }

func (le leafEntry) Write(data []byte) {
	data0 := le.recordID.Write(data)
	copy(data0, []byte(le.key))
}

func (le *leafEntry) Read(data []byte) {
	data0 := le.recordID.Read(data)
	le.key = key.Key(data0)
}
