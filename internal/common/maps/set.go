package maps

// Set on map with empty struct as value
type Set[T comparable] map[T]struct{}

func (s *Set[T]) Add(v T) {
	(*s)[v] = struct{}{}
}

func (s *Set[T]) Remove(v T) {
	delete(*s, v)
}

func (s *Set[T]) Has(v T) bool {
	_, has := (*s)[v]
	return has
}

func SetFromSlice[T any, E comparable](s []T, f func(T) E) Set[E] {
	return FromSlice(s, func(v T) (E, struct{}) {
		return f(v), struct{}{}
	})
}
