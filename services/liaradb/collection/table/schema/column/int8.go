package column

type Int8 struct {
	name string
}

func NewInt8(name string) Int8 {
	return Int8{
		name: name,
	}
}

func (i Int8) Name() string { return i.name }
func (i Int8) Size() int    { return 1 }
func (i Int8) Type() Type   { return TypeInt8 }

func (i Int8) Value() int8 {
	return 0
}
