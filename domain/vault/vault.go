package vault

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/item"
	"github.com/vterry/guild-project-ddd/domain/vault/valueobjects"
)

type Vault struct {
	valueobjects.VaultID
	Items      []*item.Item
	GoldAmount int
}

func NewVault() *Vault {
	vault := Vault{
		VaultID:    valueobjects.NewVaultID(uuid.New()),
		Items:      []*item.Item{},
		GoldAmount: 0,
	}
	return &vault
}
