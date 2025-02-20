package entity

import (
	"slices"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/valueobjects"
)

type Player struct {
	valueobjects.PlayerID
	nickname     string
	class        valueobjects.Class
	items        []*valueobjects.ItemID
	gold         int
	cash         int
	currentGuild *valueobjects.GuildID
}

func NewPlayer(nickname string, class valueobjects.Class) *Player {
	return &Player{
		PlayerID:     valueobjects.NewPlayerID(uuid.New()),
		nickname:     nickname,
		class:        class,
		items:        make([]*valueobjects.ItemID, 10),
		gold:         0,
		cash:         0,
		currentGuild: nil,
	}
}

// TODO - Improve this -- Add better error handling
func (p *Player) PickItemToInventory(item *valueobjects.ItemID) {
	p.items = append(p.items, item)
}

// TODO - Improve this -- Add better error handling
func (p *Player) RemoveItemFromInventory(item *valueobjects.ItemID) {
	itemIndex := slices.Index(p.items, item)
	p.items = append(p.items[:itemIndex], p.items[itemIndex+1:]...)
}

func (p *Player) UpdateCash(cash int) {
	p.cash = cash
}

func (p *Player) UpdateGold(gold int) {
	p.gold = gold
}

func (p *Player) UpdateCurrentGuild(g *valueobjects.GuildID) *Player {
	p.currentGuild = g
	return p
}

func (p *Player) GetCurrentGuild() uuid.UUID {
	if p.currentGuild == nil {
		return uuid.Nil
	}
	return p.currentGuild.ID()
}
