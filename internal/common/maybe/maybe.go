package maybe

type Maybe[T any] struct {
	v     T
	valid bool
}

func NewJust[T any](v T) Maybe[T] {
	return Maybe[T]{
		v:     v,
		valid: true,
	}
}

// NewNone used for explicit none value
func NewNone[T any]() Maybe[T] {
	return Maybe[T]{}
}

func Valid[T any](maybe Maybe[T]) bool {
	return maybe.valid
}

func Just[T any](maybe Maybe[T]) T {
	if !Valid(maybe) {
		panic("violated usage of maybe: Just on non Valid Maybe")
	}
	return maybe.v
}

// MapNone returns underlying value on Valid Maybe or value from f
func MapNone[T any](m Maybe[T], f func() T) T {
	if !Valid(m) {
		return f()
	}
	return Just(m)
}

func FromPtr[T any](t *T) Maybe[T] {
	if t == nil {
		return NewNone[T]()
	}
	return NewJust[T](*t)
}

func ToPtr[T any](m Maybe[T]) *T {
	if m.valid {
		return &m.v
	}
	return nil
}

func Map[T any, E any](m Maybe[T], f func(T) E) Maybe[E] {
	if !Valid(m) {
		return Maybe[E]{}
	}

	return Maybe[E]{
		v:     f(Just(m)),
		valid: true,
	}
}

func MapErr[T any, E any](m Maybe[T], f func(T) (E, error)) (Maybe[E], error) {
	if !Valid(m) {
		return Maybe[E]{}, nil
	}

	e, err := f(Just(m))
	if err != nil {
		return Maybe[E]{}, err
	}

	return Maybe[E]{
		v:     e,
		valid: true,
	}, nil
}
