package vault

import (
	"sync"

	"github.com/vterry/guild-project-ddd/domain/item"
	"github.com/vterry/guild-project-ddd/domain/player"
	"github.com/vterry/guild-project-ddd/domain/vault/valueobjects"
)

type Vault struct {
	valueobjects.VaultID
	Items      []*item.Item
	GoldAmount int
	sync.Mutex
}

func NewVault(guildName string) *Vault {
	vault := Vault{
		VaultID:    valueobjects.NewVaultID(guildName),
		Items:      []*item.Item{},
		GoldAmount: 0,
	}
	return &vault
}

func (v *Vault) AddItem(i *item.Item, p *player.Player) error {
	v.Lock()
	defer v.Unlock()

	err := p.RetriveItem(i)

	if err != nil {
		return NewVaultError(ErrInvalidOperation, err)
	}

	v.Items = append(v.Items, i)

	return nil
}

func (v *Vault) RetriveItem(i *item.Item, p *player.Player) error {
	v.Lock()
	defer v.Unlock()

	if len(v.Items) == 0 {
		return NewVaultError(ErrEmptyVault, nil)
	}

	for index, item := range v.Items {
		if item.Equals(i.ItemID) {
			if err := p.PickItem(i); err != nil {
				return NewVaultError(ErrInvalidOperation, err)
			}
			v.Items = append(v.Items[:index], v.Items[index+1:]...)
			return nil
		}
	}
	return NewVaultError(ErrItemNotFound, nil)
}

func (v *Vault) AddGold(goldAmount int, p *player.Player) error {
	v.Lock()
	defer v.Unlock()

	if goldAmount < 0 {
		return NewVaultError(ErrNegativeGoldAmount, nil)
	}

	if p.GetCurrentGold() < goldAmount {
		return NewVaultError(ErrInvalidGoldAmount, nil)
	}

	v.GoldAmount += goldAmount
	p.UpdateGold(-goldAmount)

	return nil
}

func (v *Vault) GoldWithdraw(goldAmount int, p *player.Player) error {
	v.Lock()
	defer v.Unlock()

	if goldAmount < 0 {
		return NewVaultError(ErrNegativeGoldAmount, nil)
	}

	if v.GoldAmount < goldAmount {
		return NewVaultError(ErrInvalidGoldAmount, nil)
	}

	v.GoldAmount -= goldAmount
	p.UpdateGold(goldAmount)

	return nil
}
