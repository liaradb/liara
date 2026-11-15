package column

type UInt32 struct {
	name string
}

func NewUInt32(name string) UInt32 {
	return UInt32{
		name: name,
	}
}

func (u UInt32) Name() string { return u.name }
func (u UInt32) Size() int    { return 4 }
func (u UInt32) Type() Type   { return TypeUInt32 }

func (u UInt32) Value() uint32 {
	return 0
}
