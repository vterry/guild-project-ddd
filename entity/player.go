package entity

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/valueobjects"
)

type Player struct {
	valueobjects.PlayerID
	Nickname     string
	Class        valueobjects.Class
	CurrentGuild *valueobjects.GuildID
}

func NewPlayer(nickname string, class valueobjects.Class) *Player {
	return &Player{
		PlayerID:     valueobjects.NewPlayerID(uuid.New()),
		Nickname:     nickname,
		Class:        class,
		CurrentGuild: nil,
	}
}

func (p *Player) GetCurrentGuild() string {
	if p.CurrentGuild == nil {
		return "Não pertence a nenhuma guild"
	}
	return p.CurrentGuild.ID().String()
}
