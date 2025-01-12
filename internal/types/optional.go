package types

type Optional[T any] struct {
	Value    T
	HasValue bool
}

func ToOptional[T any](value T) Optional[T] {
	return Optional[T]{Value: value, HasValue: true}
}
