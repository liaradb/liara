package storage

const (
	Uint64Length = 8
)

type Uint64Entry struct {
	Position Position
}

func newUInt64Entry(p Position) Uint64Entry {
	return Uint64Entry{Position: p}
}

func (u Uint64Entry) Next() Position {
	return u.Position + Uint64Length
}

func (u Uint64Entry) Get(b *Buffer) (uint64, error) {
	return b.ReadUint64(int(u.Position))
}

func (u Uint64Entry) Set(b *Buffer, value uint64) error {
	return b.WriteUint64(value, int(u.Position))
}
