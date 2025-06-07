package base

import (
	"crypto/rand"
	"encoding/base64"
)

type BaseID[T comparable] struct {
	value T
}

func New[T comparable](value T) BaseID[T] {
	return BaseID[T]{value: value}
}

func (b BaseID[T]) ID() T {
	return b.value
}

func (this BaseID[T]) Equals(object any) bool {
	o, ok := object.(BaseID[T])
	return ok && this.value == o.value
}

func ShortUUID(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b)[:length*4/3], nil
}
