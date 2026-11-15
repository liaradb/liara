package column

type String struct {
	name string
	size int
}

func NewString(name string, size int) String {
	return String{
		name: name,
		size: size,
	}
}

func (s String) Name() string { return s.name }
func (s String) Size() int    { return s.size }
func (s String) Type() Type   { return TypeString }

func (s String) Value() string {
	return ""
}
