package key

import "fmt"

type intKey int64

func (k intKey) Size() int      { return 8 }
func (k intKey) String() string { return fmt.Sprintf("%v", rune(k)) }
func (k intKey) Value() int64   { return int64(k) }

func (k intKey) Equal(o any) bool {
	b, ok := o.(intKey)
	if !ok {
		return false
	}

	return k == b
}

func (k intKey) Greater(o any) bool {
	b, ok := o.(intKey)
	if !ok {
		return false
	}

	return k > b
}

func (k intKey) GreaterEqual(o any) bool {
	b, ok := o.(intKey)
	if !ok {
		return false
	}

	return k >= b
}

func (k intKey) Less(o any) bool {
	b, ok := o.(intKey)
	if !ok {
		return false
	}

	return k < b
}

func (k intKey) LessEqual(o any) bool {
	b, ok := o.(intKey)
	if !ok {
		return false
	}

	return k <= b
}
