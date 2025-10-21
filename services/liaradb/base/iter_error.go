package base

import "iter"

func IterError[T any](err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var t T
		yield(t, err)
	}
}
