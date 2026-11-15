package column

type Int64 struct {
	name string
}

func NewInt64(name string) Int64 {
	return Int64{
		name: name,
	}
}

func (i Int64) Name() string { return i.name }
func (i Int64) Size() int    { return 8 }
func (i Int64) Type() Type   { return TypeInt64 }

func (i Int64) Value() int64 {
	return 0
}
