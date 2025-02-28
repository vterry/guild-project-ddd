package player

import (
	"slices"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
	item "github.com/vterry/guild-project-ddd/domain/item/valueobjects"
	"github.com/vterry/guild-project-ddd/domain/player/valueobjects"
)

type Player struct {
	valueobjects.PlayerID
	nickname     string
	class        common.Class
	items        []*item.ItemID
	gold         int
	cash         int
	currentGuild uuid.UUID
}

func NewPlayer(nickname string, class common.Class) *Player {
	return &Player{
		PlayerID:     valueobjects.NewPlayerID(uuid.New()),
		nickname:     nickname,
		class:        class,
		items:        make([]*item.ItemID, 10),
		gold:         0,
		cash:         0,
		currentGuild: uuid.Nil,
	}
}

// TODO - Improve this -- Add better error handling
func (p *Player) PickItemToInventory(item *item.ItemID) {
	p.items = append(p.items, item)
}

// TODO - Improve this -- Add better error handling
func (p *Player) RemoveItemFromInventory(item *item.ItemID) {
	itemIndex := slices.Index(p.items, item)
	p.items = append(p.items[:itemIndex], p.items[itemIndex+1:]...)
}

func (p *Player) UpdateCash(cash int) {
	p.cash = cash
}

func (p *Player) UpdateGold(gold int) {
	p.gold = gold
}

func (p *Player) UpdateCurrentGuild(g uuid.UUID) *Player {
	p.currentGuild = g
	return p
}

func (p *Player) GetCurrentGuild() uuid.UUID {
	if p.currentGuild == uuid.Nil {
		return uuid.Nil
	}
	return p.currentGuild
}

func (p *Player) GetCurrentCash() int {
	return p.cash
}
