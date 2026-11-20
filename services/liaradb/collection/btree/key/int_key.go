package key

import "fmt"

type intKey int64

func (k intKey) Size() int      { return 8 }
func (k intKey) String() string { return fmt.Sprintf("%v", rune(k)) }
func (k intKey) Value() int64   { return int64(k) }

func (k intKey) Equal(b intKey) bool {
	return k == b
}

func (k intKey) Greater(b intKey) bool {
	return k > b
}

func (k intKey) GreaterEqual(b intKey) bool {
	return k >= b
}

func (k intKey) Less(b intKey) bool {
	return k < b
}

func (k intKey) LessEqual(b intKey) bool {
	return k <= b
}
