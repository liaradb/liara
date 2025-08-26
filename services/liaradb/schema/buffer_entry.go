package storage

import (
	"github.com/cardboardrobots/liaradb/raw"
	"github.com/cardboardrobots/liaradb/storage"
)

type Uint64Entry struct {
	Offset raw.Offset
}

func newUInt64Entry(o raw.Offset) Uint64Entry {
	return Uint64Entry{Offset: o}
}

func (u Uint64Entry) Next() raw.Offset {
	return u.Offset + raw.Uint64Length
}

func (u Uint64Entry) Get(b *storage.Buffer) (uint64, error) {
	return b.ReadUint64(u.Offset)
}

func (u Uint64Entry) Set(b *storage.Buffer, value uint64) error {
	return b.WriteUint64(value, u.Offset)
}
