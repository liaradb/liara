package column

type UInt16 struct {
	name string
}

func NewUInt16(name string) UInt16 {
	return UInt16{
		name: name,
	}
}

func (u UInt16) Name() string { return u.name }
func (u UInt16) Size() int    { return 2 }
func (u UInt16) Type() Type   { return TypeUInt16 }

func (u UInt16) Value() uint16 {
	return 0
}
