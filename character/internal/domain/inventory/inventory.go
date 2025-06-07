package inventory

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/domain/playeritem"
)

const MAX_ITEMS = 10

var (
	ErrInventoryIsFull    = errors.New("cannot add item - inventory is full")
	ErrPlayerItemNotFound = errors.New("cannot drop item - player item is not in inventory")
	ErrInvalidGoldAmount  = errors.New("cannot add gold - invalid format")
	ErrNotEnoughGold      = errors.New("cannot withdraw that amount - gold is not enough")
)

type Inventory struct {
	InventoryID
	goldAmount int
	items      map[playeritem.PlayerItemID]playeritem.PlayerItem
}

func NewInventory() *Inventory {
	return &Inventory{
		InventoryID: New(uuid.New()),
		goldAmount:  0,
		items:       make(map[playeritem.PlayerItemID]playeritem.PlayerItem, MAX_ITEMS),
	}
}

func (i *Inventory) AddItem(playeritem playeritem.PlayerItem) error {

	if len(i.items)+1 > MAX_ITEMS {
		return ErrInventoryIsFull
	}

	i.items[playeritem.PlayerItemID] = playeritem

	return nil
}

func (i *Inventory) DropItem(playeritem playeritem.PlayerItem) error {

	if _, ok := i.items[playeritem.PlayerItemID]; !ok {
		return ErrPlayerItemNotFound
	}

	delete(i.items, playeritem.PlayerItemID)

	return nil
}

func (i *Inventory) AddGold(amount int) error {
	if amount < 0 {
		return ErrInvalidGoldAmount
	}

	i.goldAmount += amount
	return nil
}

func (i *Inventory) WithdrawGold(amount int) error {

	if amount < 0 {
		return ErrInvalidGoldAmount
	}

	if i.goldAmount-amount < 0 {
		return ErrNotEnoughGold
	}

	i.goldAmount -= amount
	return nil
}

func (i *Inventory) GetCurrentGold() int {
	return i.goldAmount
}

func (i *Inventory) ShowItems() []playeritem.PlayerItemID {
	keys := make([]playeritem.PlayerItemID, 0, len(i.items))
	for k := range i.items {
		keys = append(keys, k)
	}
	return keys
}
