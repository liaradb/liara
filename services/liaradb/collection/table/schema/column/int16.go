package column

type Int16 struct {
	name string
}

func NewInt16(name string) Int16 {
	return Int16{
		name: name,
	}
}

func (i Int16) Name() string { return i.name }
func (i Int16) Size() int    { return 2 }
func (i Int16) Type() Type   { return TypeInt16 }

func (i Int16) Value() int16 {
	return 0
}
