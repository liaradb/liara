package column

type Int32 struct {
	name string
}

func NewInt32(name string) Int32 {
	return Int32{
		name: name,
	}
}

func (i Int32) Name() string { return i.name }
func (i Int32) Size() int    { return 4 }
func (i Int32) Type() Type   { return TypeInt32 }

func (i Int32) Value() int32 {
	return 0
}
