package wrap

type Bool byte

func (b *Bool) Set(i byte, value bool) {
	if value {
		*b |= 1 << i
	} else {
		*b &^= 1 << i
	}
}

func (b Bool) Get(i byte) bool {
	return (b>>i)&1 == 1
}
