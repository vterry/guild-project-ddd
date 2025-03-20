package player

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/item"
	"github.com/vterry/guild-project-ddd/domain/player/valueobjects"
)

const MAX_ITEMS = 10

var (
	ErrInvalidNickname = errors.New("nickname must by between 4 and 15 characters")
	ErrNotEnoughSpace  = errors.New("player has no space for more items")
	ErrItemNotFound    = errors.New("item not found in player inventory")
	ErrEmptyInventory  = errors.New("player has no items")
	ErrNotEnoughGold   = errors.New("player has not enough gold")
	ErrNotEnoughCash   = errors.New("player has not enough cash")
)

type Player struct {
	valueobjects.PlayerID
	nickname     string
	class        valueobjects.Class
	items        []*item.Item
	gold         int
	cash         int
	currentGuild string
	sync.Mutex
}

func NewPlayer(nickname string, class valueobjects.Class) (*Player, error) {
	if !isValidNickname(nickname) {
		return nil, ErrInvalidNickname
	}

	player := initializePlayer(nickname, class)
	return player, nil
}

func (p *Player) GetCurrentGuild() string {
	return p.currentGuild
}

func (p *Player) UpdateCurrentGuild(g string) *Player {
	p.Lock()
	defer p.Unlock()
	p.currentGuild = g
	return p
}

func (p *Player) UpdateCash(cash int) error {
	p.Lock()
	defer p.Unlock()

	nxtAmount := p.cash + cash

	if nxtAmount <= 0 {
		return ErrNotEnoughCash
	}

	p.cash = nxtAmount
	return nil
}

func (p *Player) UpdateGold(gold int) error {
	p.Lock()
	defer p.Unlock()
	nxtAmount := p.gold + gold

	if nxtAmount <= 0 {
		return ErrNotEnoughGold
	}

	p.gold = nxtAmount
	return nil
}

func (p *Player) GetCurrentCash() int {
	return p.cash
}

func (p *Player) GetCurrentGold() int {
	return p.gold
}

func (p *Player) PickItem(i *item.Item) error {
	p.Lock()
	defer p.Unlock()

	if len(p.items) >= MAX_ITEMS {
		return ErrNotEnoughSpace
	}
	p.items = append(p.items, i)
	return nil
}

func (p *Player) RetriveItem(i *item.Item) error {
	p.Lock()
	defer p.Unlock()

	if len(p.items) == 0 {
		return ErrEmptyInventory
	}

	for index, item := range p.items {
		if item == i {
			p.items = append(p.items[:index], p.items[index+1:]...)
			break
		} else {
			return ErrItemNotFound
		}
	}
	return nil
}

func isValidNickname(nickname string) bool {
	return len(nickname) > 0 && len(nickname) < 15
}

func initializePlayer(nickname string, class valueobjects.Class) *Player {
	return &Player{
		PlayerID:     valueobjects.NewPlayerID(uuid.New()),
		nickname:     nickname,
		class:        class,
		items:        make([]*item.Item, 0, MAX_ITEMS),
		gold:         0,
		cash:         0,
		currentGuild: "",
	}
}
