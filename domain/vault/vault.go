package vault

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/item"
	"github.com/vterry/guild-project-ddd/domain/player"
	"github.com/vterry/guild-project-ddd/domain/vault/valueobjects"
)

var (
	ErrInvalidGoldAmount = errors.New("player does not have the current gold amount")
)

type Vault struct {
	valueobjects.VaultID
	Items      []*item.Item
	GoldAmount int
	sync.Mutex
}

func NewVault() *Vault {
	vault := Vault{
		VaultID:    valueobjects.NewVaultID(uuid.New()),
		Items:      []*item.Item{},
		GoldAmount: 0,
	}
	return &vault
}

func (v *Vault) AddItem(i *item.Item, p *player.Player) {
	v.Lock()
	defer v.Unlock()
	v.Items = append(v.Items, i)
	p.RetrieItem(i)

}

func (v *Vault) RetrieItem(i *item.Item, p *player.Player) error {
	v.Lock()
	defer v.Unlock()

	for index, item := range v.Items {
		if item.ID() == i.ID() {
			v.Items = append(v.Items[:index], v.Items[index+1:]...)
			if err := p.PickItem(i); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("item not found in vault")
}

func (v *Vault) AddGold(goldAmount int, p *player.Player) error {
	v.Lock()
	defer v.Unlock()

	if p.GetCurrentGold() < goldAmount {
		return ErrInvalidGoldAmount
	}

	v.GoldAmount += goldAmount
	p.UpdateGold(p.GetCurrentGold() - goldAmount)

	return nil
}

func (v *Vault) GoldWithdraw(goldAmount int, p *player.Player) error {
	v.Lock()
	defer v.Unlock()

	if v.GoldAmount < goldAmount {
		return ErrInvalidGoldAmount
	}

	v.GoldAmount -= goldAmount
	p.UpdateGold(p.GetCurrentGold() + goldAmount)

	return nil
}
