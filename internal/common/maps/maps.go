package maps

func FromSlice[K comparable, V, E any](s []V, f func(V) (K, E)) map[K]E {
	res := map[K]E{}
	for _, v := range s {
		k, newV := f(v)
		res[k] = newV
	}
	return res
}

func ToSlice[K comparable, V, E any](m map[K]E, f func(K, E) V) []V {
	res := make([]V, 0, len(m))
	for k, e := range m {
		res = append(res, f(k, e))
	}
	return res
}
