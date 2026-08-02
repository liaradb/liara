package span

import (
	"io"

	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/encoder/raw"
)

type testRecord struct {
	id   base.Uint64
	data []byte
}

func (tr *testRecord) ID() base.Uint64 {
	return tr.id
}

func (tr *testRecord) Read(r io.Reader) error {
	if err := tr.id.Read(r); err != nil {
		return err
	}

	return raw.Read(r, &tr.data)
}

func (tr *testRecord) SetLogSequenceNumber(id base.Uint64) {
	tr.id = id
}

func (tr *testRecord) Size() int {
	return base.Uint64Size +
		raw.HeaderSize + len(tr.data)
}

func (tr *testRecord) Write(w io.Writer) error {
	if err := tr.id.Write(w); err != nil {
		return err
	}

	return raw.Write(w, tr.data)
}
