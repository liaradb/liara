package base

import "iter"

func IterToSlice[T any](i iter.Seq2[T, error]) ([]T, int, error) {
	s := make([]T, 0)
	for t, err := range i {
		if err != nil {
			return nil, 0, err
		}

		s = append(s, t)
	}
	return s, 0, nil
}
