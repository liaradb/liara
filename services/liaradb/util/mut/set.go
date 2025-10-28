package mut

type (
	Set[T comparable] map[T]struct{}
)

func NewSet[T comparable](t ...T) Set[T] {
	s := Set[T]{}
	s.Add(t...)
	return s
}

// Must not be nil
func (s Set[T]) Add(t ...T) {
	for _, i := range t {
		s[i] = struct{}{}
	}
}

func (s Set[T]) Remove(t ...T) {
	for _, i := range t {
		delete(s, i)
	}
}

func (s Set[T]) Clear() {
	clear(s)
}

func (s Set[T]) Includes(t T) bool {
	_, ok := s[t]
	return ok
}

func (s Set[T]) Slice() []T {
	values := make([]T, 0, len(s))
	for t := range s {
		values = append(values, t)
	}
	return values
}

// Must not be nil
func (s Set[T]) Merge(b Set[T]) {
	for i := range b {
		s[i] = struct{}{}
	}
}

func (s Set[T]) Union(b Set[T]) Set[T] {
	l := len(s)
	if lb := len(b); lb > l {
		l = lb
	}
	c := make(Set[T], l)

	for i := range s {
		c[i] = struct{}{}
	}

	for i := range b {
		c[i] = struct{}{}
	}

	return c
}

func (s Set[T]) Intersection(b Set[T]) Set[T] {
	l := len(s)
	if lb := len(b); lb < l {
		l = lb
	}
	c := make(Set[T], l)

	for i := range s {
		if _, ok := b[i]; ok {
			c[i] = struct{}{}
		}
	}

	return c
}
