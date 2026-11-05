package column

type Byte struct {
	name string
}

func NewByte(name string) Byte {
	return Byte{
		name: name,
	}
}

func (b Byte) Name() string { return b.name }
func (b Byte) Size() int    { return 1 }
func (b Byte) Type() Type   { return TypeByte }

func (b Byte) Value() byte {
	return 0
}
