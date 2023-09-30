package slices

import (
	"github.com/ispringtech/brewkit/internal/common/maybe"
)

// Map iterates through slice and maps values
func Map[T, TResult any](s []T, f func(T) TResult) []TResult {
	result := make([]TResult, 0, len(s))
	for _, t := range s {
		result = append(result, f(t))
	}
	return result
}

// MapErr iterates through slice and maps values and stops on any error
func MapErr[T, TResult any](s []T, f func(T) (TResult, error)) ([]TResult, error) {
	result := make([]TResult, 0, len(s))
	for _, t := range s {
		e, err := f(t)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

// Filter iterates and adds to result slice elements that satisfied predicate
func Filter[T any](s []T, f func(T) bool) []T {
	var result []T
	for _, t := range s {
		if f(t) {
			result = append(result, t)
		}
	}
	return result
}

// FilterErr iterates and adds to result slice elements that satisfied predicate and stop on any error
func FilterErr[T any](s []T, f func(T) (bool, error)) ([]T, error) {
	var result []T
	for _, t := range s {
		accepted, err := f(t)
		if err != nil {
			return nil, err
		}

		if accepted {
			result = append(result, t)
		}
	}

	return result, nil
}

func MapMaybe[T, TResult any](s []T, f func(T) maybe.Maybe[TResult]) (res []TResult) {
	for _, t := range s {
		m := f(t)
		if maybe.Valid(m) {
			res = append(res, maybe.Just(m))
		}
	}
	return res
}

func MapMaybeErr[T, TResult any](s []T, f func(T) (maybe.Maybe[TResult], error)) (res []TResult, err error) {
	for _, t := range s {
		var m maybe.Maybe[TResult]
		m, err = f(t)
		if err != nil {
			return nil, err
		}
		if maybe.Valid(m) {
			res = append(res, maybe.Just(m))
		}
	}
	return res, nil
}

// Merge slices into one
func Merge[T any](slices ...[]T) []T {
	if len(slices) == 0 {
		return nil
	}

	result := slices[0]
	for i := 1; i < len(slices); i++ {
		s := slices[i]

		result = append(result, s...)
	}

	return result
}
