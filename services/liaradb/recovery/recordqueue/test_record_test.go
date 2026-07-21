package recordqueue

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/recovery/logpage"
)

type testRecord struct {
	lsn  logpage.LogSequenceNumber
	data []byte
}

func newTestRecord(size int) *testRecord {
	return &testRecord{
		data: make([]byte, size),
	}
}

func newTestRecordData(data []byte) *testRecord {
	return &testRecord{
		data: data,
	}
}

func (tr *testRecord) LogSequenceNumber() logpage.LogSequenceNumber {
	return tr.lsn
}

func (tr *testRecord) Read(r io.Reader) error {
	if err := tr.lsn.Read(r); err != nil {
		return err
	}

	return raw.Read(r, &tr.data)
}

func (tr *testRecord) SetLogSequenceNumber(lsn logpage.LogSequenceNumber) {
	tr.lsn = lsn
}

func (tr *testRecord) Size() int {
	return logpage.LogSequenceNumberSize +
		raw.HeaderSize + len(tr.data)
}

func (tr *testRecord) Write(w io.Writer) error {
	if err := tr.lsn.Write(w); err != nil {
		return err
	}

	return raw.Write(w, tr.data)
}
