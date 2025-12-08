package btree

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

func (le LeafEntry) Write(data []byte) {
	data0 := le.recordID.Write(data)
	copy(data0, []byte(le.key))
}

func (le *LeafEntry) Read(data []byte) {
	data0 := le.recordID.Read(data)
	le.key = Key(data0)
}
