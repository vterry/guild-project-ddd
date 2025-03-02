package vault

import (
	"errors"

	"github.com/vterry/guild-project-ddd/domain/common"
)

var (
	ErrInvalidGoldAmount  = errors.New("player does not have the current gold amount")
	ErrNegativeGoldAmount = errors.New("gold amount cannot be negative")
	ErrInvalidOperation   = errors.New("invalid operation")
	ErrItemNotFound       = errors.New("item not found in vault")
	ErrEmptyVault         = errors.New("vault's item is empty")
)

type VaultError struct {
	common.BaseError
}

func NewVaultError(curErr, srcErr error) *VaultError {
	return &VaultError{
		BaseError: common.NewError(curErr, srcErr).(common.BaseError),
	}
}
