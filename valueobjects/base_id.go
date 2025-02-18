package valueobjects

type BaseID[T comparable] struct {
	value T
}

func NewBaseID[T comparable](value T) BaseID[T] {
	return BaseID[T]{value: value}
}

func (b BaseID[T]) ID() T {
	return b.value
}

func (this BaseID[T]) Equals(object any) bool {
	o, ok := object.(BaseID[T])
	return ok && this.value == o.value
}
