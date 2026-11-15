package column

type UInt64 struct {
	name string
}

func NewUInt64(name string) UInt64 {
	return UInt64{
		name: name,
	}
}

func (u UInt64) Name() string { return u.name }
func (u UInt64) Size() int    { return 8 }
func (u UInt64) Type() Type   { return TypeUInt64 }

func (u UInt64) Value() uint64 {
	return 0
}
