package common

import "errors"

type BaseError struct {
	curErr error
	srcErr error
}

func NewError(curErr, srcErr error) error {
	return BaseError{
		srcErr: srcErr,
		curErr: curErr,
	}
}

func (e BaseError) Error() string {
	return errors.Join(e.curErr, e.srcErr).Error()
}
