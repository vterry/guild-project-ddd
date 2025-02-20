package entity

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/valueobjects"
)

type Vault struct {
	valueobjects.VaultID
	Items      []*Item
	GoldAmount int
}

func NewVault() *Vault {
	vault := Vault{
		VaultID:    valueobjects.NewVaultID(uuid.New()),
		Items:      []*Item{},
		GoldAmount: 0,
	}
	return &vault
}
