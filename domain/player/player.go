package player

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
	"github.com/vterry/guild-project-ddd/domain/item"
	"github.com/vterry/guild-project-ddd/domain/player/valueobjects"
)

var (
	ErrInvalidNickname = errors.New("nickname must by between 4 and 15 characters")
	ErrNotEnoughSpace  = errors.New("player has no space for more items")
	MAX_ITEMS          = 10
)

type Player struct {
	valueobjects.PlayerID
	nickname     string
	class        common.Class
	items        []*item.Item
	gold         *int
	cash         *int
	currentGuild uuid.UUID
	sync.Mutex
}

func NewPlayer(nickname string, class common.Class) (*Player, error) {
	if !isValidNickname(nickname) {
		return nil, ErrInvalidNickname
	}

	player := initializePlayer(nickname, class)
	return player, nil
}

func (p *Player) GetCurrentGuild() uuid.UUID {
	return p.currentGuild
}

func (p *Player) UpdateCurrentGuild(g uuid.UUID) *Player {
	p.Lock()
	defer p.Unlock()

	p.currentGuild = g
	return p
}

func (p *Player) UpdateCash(cash int) {
	p.Lock()
	defer p.Unlock()

	p.cash = &cash
}

func (p *Player) UpdateGold(gold int) {
	p.Lock()
	defer p.Unlock()

	p.gold = &gold
}

func (p *Player) GetCurrentCash() int {
	return *p.cash
}

func (p *Player) GetCurrentGold() int {
	return *p.gold
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

func (p *Player) RetrieItem(i *item.Item) {
	p.Lock()
	defer p.Unlock()

	for index, item := range p.items {
		if item == i {
			p.items = append(p.items[:index], p.items[index+1:]...)
		}
	}
}

func isValidNickname(nickname string) bool {
	return len(nickname) > 0 && len(nickname) < 15
}

func initializePlayer(nickname string, class common.Class) *Player {
	return &Player{
		PlayerID:     valueobjects.NewPlayerID(uuid.New()),
		nickname:     nickname,
		class:        class,
		items:        make([]*item.Item, MAX_ITEMS),
		gold:         new(int),
		cash:         new(int),
		currentGuild: uuid.Nil,
	}
}
