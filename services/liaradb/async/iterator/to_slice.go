package iterator

import "iter"

func ToSlice[T any](i iter.Seq2[T, error]) ([]T, error) {
	s := make([]T, 0)
	for t, err := range i {
		if err != nil {
			return nil, err
		}

		s = append(s, t)
	}
	return s, nil
}
