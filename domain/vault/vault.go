package vault

import (
	"errors"
	"fmt"
	"sync"

	"github.com/vterry/guild-project-ddd/domain/item"
	"github.com/vterry/guild-project-ddd/domain/player"
	"github.com/vterry/guild-project-ddd/domain/vault/valueobjects"
)

var (
	ErrInvalidGoldAmount  = errors.New("player does not have the current gold amount")
	ErrNegativeGoldAmount = errors.New("gold amount cannot be negative")
	ErrInvalidOperation   = errors.New("invalid operation")
	ErrItemNotFound       = errors.New("item not found in vault")
	ErrEmptyVault         = errors.New("vault's item is empty")
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
		return fmt.Errorf("%w: %v", ErrInvalidOperation, err)
	}

	v.Items = append(v.Items, i)

	return nil
}

func (v *Vault) RetriveItem(i *item.Item, p *player.Player) error {
	v.Lock()
	defer v.Unlock()

	if len(v.Items) == 0 {
		return ErrEmptyVault
	}

	for index, item := range v.Items {
		if item.Equals(i.ItemID) {
			if err := p.PickItem(i); err != nil {
				return fmt.Errorf("%w: %v", ErrInvalidOperation, err)
			}
			v.Items = append(v.Items[:index], v.Items[index+1:]...)
			return nil
		}
	}
	return ErrItemNotFound
}

func (v *Vault) AddGold(goldAmount int, p *player.Player) error {
	v.Lock()
	defer v.Unlock()

	if goldAmount < 0 {
		return ErrNegativeGoldAmount
	}

	if p.GetCurrentGold() < goldAmount {
		return ErrInvalidGoldAmount
	}

	v.GoldAmount += goldAmount
	p.UpdateGold(-goldAmount)

	return nil
}

func (v *Vault) GoldWithdraw(goldAmount int, p *player.Player) error {
	v.Lock()
	defer v.Unlock()

	if goldAmount < 0 {
		return ErrNegativeGoldAmount
	}

	if v.GoldAmount < goldAmount {
		return ErrInvalidGoldAmount
	}

	v.GoldAmount -= goldAmount
	p.UpdateGold(goldAmount)

	return nil
}
