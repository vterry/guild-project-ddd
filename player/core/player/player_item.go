package player

import "github.com/vterry/ddd-study/player/core/common/valueobjects"

type PlayerItem struct {
	ItemId   valueobjects.ItemID
	quantity int
}

func NewPlayerItem(itemId valueobjects.ItemID, quantity int) *PlayerItem {
	return &PlayerItem{
		ItemId:   itemId,
		quantity: quantity,
	}
}
