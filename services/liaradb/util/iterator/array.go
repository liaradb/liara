package iterator

import "fmt"

func Map[T any, U any](a []T, p func(a T) U) []U {
	length := len(a)
	if length == 0 {
		return nil
	}

	b := make([]U, 0, length)
	for _, i := range a {
		b = append(b, p(i))
	}

	return b
}

func MapMap[T any, U any](a map[any]T, p func(a T) U) []U {
	length := len(a)
	if length == 0 {
		return nil
	}

	b := make([]U, 0, length)
	for _, i := range a {
		b = append(b, p(i))
	}

	return b
}

func MapValues[T any, U comparable](a map[U]T) []T {
	length := len(a)
	if length == 0 {
		return nil
	}

	b := make([]T, 0, length)
	for _, i := range a {
		b = append(b, i)
	}

	return b
}

type (
	HashMap[K comparable, V any] map[K]V
)

func (m HashMap[K, V]) Values() []V {
	length := len(m)
	if length == 0 {
		return nil
	}

	v := make([]V, 0, length)
	for _, i := range m {
		v = append(v, i)
	}

	return v
}

func SliceToMap[T comparable, U any](slice []U, key func(U) T) map[T]U {
	result := make(map[T]U, len(slice))
	for _, value := range slice {
		result[key(value)] = value
	}
	return result
}

func MapString[T fmt.Stringer](a []T) []string {
	length := len(a)
	if length == 0 {
		return nil
	}

	b := make([]string, 0, length)
	for _, i := range a {
		b = append(b, i.String())
	}

	return b
}

func MapFromString[T ~string](a []string) []T {
	length := len(a)
	if length == 0 {
		return nil
	}

	b := make([]T, 0, length)
	for _, i := range a {
		b = append(b, T(i))
	}

	return b
}

func MapToString[T ~string](a []T) []string {
	length := len(a)
	if length == 0 {
		return nil
	}

	b := make([]string, 0, length)
	for _, i := range a {
		b = append(b, string(i))
	}

	return b
}
