package leafnode

import "github.com/liaradb/liaradb/collection/btree/key"

// TODO: Test this
// TODO: This should be private
type LeafEntry struct {
	key      key.Key
	recordID RecordID
}

func newLeafEntry(key key.Key, recordID RecordID) LeafEntry {
	return LeafEntry{
		key:      key,
		recordID: recordID,
	}
}

func (le LeafEntry) Key() key.Key       { return le.key }
func (le LeafEntry) RecordID() RecordID { return le.recordID }
func (le LeafEntry) Size() int          { return le.key.Size() + RecordIDSize }

func (le LeafEntry) Write(data []byte) {
	data0 := le.recordID.Write(data)
	copy(data0, []byte(le.key))
}

func (le *LeafEntry) Read(data []byte) {
	data0 := le.recordID.Read(data)
	le.key = key.Key(data0)
}
