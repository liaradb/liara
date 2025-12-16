package key

type CompositeKey struct {
	A Key
	B Key
}

func (k CompositeKey) Length() int16 {
	return k.A.Length() + k.B.Length()
}

func (k CompositeKey) Equal(o any) bool {
	b, ok := o.(CompositeKey)
	if !ok {
		return false
	}

	return k.A.Equal(b.A) && k.B.Equal(b.B)
}

func (k CompositeKey) Greater(o any) bool {
	b, ok := o.(CompositeKey)
	if !ok {
		return false
	}

	return k.A.Greater(b.A) && k.B.Greater(b.B)
}

func (k CompositeKey) GreaterEqual(o any) bool {
	b, ok := o.(CompositeKey)
	if !ok {
		return false
	}

	return k.A.GreaterEqual(b.A) && k.B.GreaterEqual(b.B)
}

func (k CompositeKey) Less(o any) bool {
	b, ok := o.(CompositeKey)
	if !ok {
		return false
	}

	return k.A.Less(b.A) && k.B.Less(b.B)
}

func (k CompositeKey) LessEqual(o any) bool {
	b, ok := o.(CompositeKey)
	if !ok {
		return false
	}

	return k.A.LessEqual(b.A) && k.B.LessEqual(b.B)
}
