package either

type discriminator uint

const (
	left discriminator = iota
	right
)

func NewLeft[L any, R any](l L) Either[L, R] {
	return Either[L, R]{
		l: l,
		d: left,
	}
}

func NewEitherLeft[T Either[L, R], L any, R any](l L) T {
	return T(Either[L, R]{
		l: l,
		d: left,
	})
}

func NewRight[L any, R any](r R) Either[L, R] {
	return Either[L, R]{
		r: r,
		d: right,
	}
}

func NewEitherRight[T Either[L, R], L any, R any](r R) T {
	return T(Either[L, R]{
		r: r,
		d: right,
	})
}

type Either[L any, R any] struct {
	l L
	r R
	d discriminator
}

func (e Either[L, R]) MapLeft(f func(l L)) Either[L, R] {
	if e.d == left {
		f(e.l)
	}
	return e
}

func (e Either[L, R]) MapRight(f func(r R)) Either[L, R] {
	if e.d == right {
		f(e.r)
	}
	return e
}

func (e Either[T, E]) IsLeft() bool {
	return e.d == left
}

func (e Either[T, E]) IsRight() bool {
	return e.d == right
}

func (e Either[T, E]) Left() T {
	if e.d != left {
		panic("violated usage of either")
	}
	return e.l
}

func (e Either[T, E]) Right() E {
	if e.d != right {
		panic("violated usage of either")
	}
	return e.r
}
