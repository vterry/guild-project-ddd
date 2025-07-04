package playeritem

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/base"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/item"
)

var (
	ErrNilItemId      = errors.New("an item id was not provided - please inform an item id")
	ErrNilDescription = errors.New("an item description must be provided - please inform an item description")
	ErrNilQuantity    = errors.New("an item quantity must be provided - please inform an item quantity")
)

type PlayerItemID struct {
	base.BaseID[uuid.UUID]
}

func NewPlayerItemID(value uuid.UUID) PlayerItemID {
	return PlayerItemID{
		BaseID: base.New(value),
	}
}

type PlayerItem struct {
	PlayerItemID
	itemID      item.ItemID
	description string
	quantity    int
}

func NewPlayerItem(itemID item.ItemID, description string, quantity int) (*PlayerItem, error) {

	if itemID.ID() == uuid.Nil {
		return nil, ErrNilItemId
	}

	if description == "" {
		return nil, ErrNilDescription
	}

	if quantity == 0 {
		return nil, ErrNilQuantity
	}

	return &PlayerItem{
		PlayerItemID: NewPlayerItemID(uuid.New()),
		itemID:       item.NewItemID(itemID.ID()),
		description:  description,
		quantity:     quantity,
	}, nil
}

func (p *PlayerItem) Describe() string {
	return p.description
}

func (p *PlayerItem) GetCurrentQuantity() int {
	return p.quantity
}

func (p *PlayerItem) ItemID() item.ItemID {
	return p.itemID
}

func (playerItemID PlayerItemID) Equals(other PlayerItemID) bool {
	return playerItemID.BaseID.Equals(other.BaseID)
}
