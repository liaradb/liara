package esmongo

type decoder interface {
	Decode(any) error
}

func decode[M any](decoder decoder) (M, error) {
	var m M
	return m, decoder.Decode(&m)
}
