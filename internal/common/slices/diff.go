package slices

// Diff return diff elements from slices
func Diff[T comparable](s1, s2 []T) []T {
	s1Map := map[T]struct{}{}
	for _, element := range s1 {
		s1Map[element] = struct{}{}
	}
	s2Map := map[T]struct{}{}
	for _, element := range s2 {
		s2Map[element] = struct{}{}
	}

	var res []T
	for _, element := range s2 {
		if _, exists := s1Map[element]; !exists {
			res = append(res, element)
		}
	}
	for _, element := range s1 {
		if _, exists := s2Map[element]; !exists {
			res = append(res, element)
		}
	}
	return res
}
