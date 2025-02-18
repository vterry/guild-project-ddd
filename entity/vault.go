package entity

import (
	"github.com/vterry/guild-project-ddd/valueobjects"
)

type Vault struct {
	valueobjects.VaultID
	Items      []*Item
	GoldAmount int
}
