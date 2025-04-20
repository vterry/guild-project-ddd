package common

import "errors"

type Base[T any] struct {
	Entity *T
}

type Specification[T any] func(Base[T]) error

func And[T any](specs ...Specification[T]) Specification[T] {
	return func(b Base[T]) error {
		var errs error
		for _, spec := range specs {
			if err := spec(b); err != nil {
				errs = errors.Join(errs, err)
			}
		}
		return errs
	}
}
