package vault

import (
	"errors"
)

var (
	ErrInvalidGoldAmount  = errors.New("player does not have the current gold amount")
	ErrNegativeGoldAmount = errors.New("gold amount cannot be negative")
	ErrItemNotFound       = errors.New("item not found in vault")
	ErrEmptyVault         = errors.New("vault's item is empty")
	ErrInvalidOperation   = errors.New("invalid operation")
)

type Vault struct {
	VaultID
	Items      []*VaultItem
	GoldAmount int
}

func newVault(guildName string) *Vault {
	vault := Vault{
		VaultID:    NewVaultID(guildName),
		Items:      []*VaultItem{},
		GoldAmount: 0,
	}
	return &vault
}
