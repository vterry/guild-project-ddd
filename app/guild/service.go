package service

import (
	"github.com/google/uuid"
	guild "github.com/vterry/guild-project-ddd/app/guild/repository"
	player "github.com/vterry/guild-project-ddd/app/player/repository"
)

type GuildService struct {
	guilds  guild.GuildRepository
	players player.PlayerRepository
}

func GuildServiceFactory(cfgs ...GuildServiceConfiguration) (*GuildService, error) {
	gs := &GuildService{}

	for _, cfgs := range cfgs {
		err := cfgs(gs)
		if err != nil {
			return nil, err
		}
	}
	return gs, nil
}

// Criar guild
func (g *GuildService) CreateGuild(guildName string, guildOwnerId uuid.UUID) {

}

// Enviar invite

// Aprovar player

// Recusar invite

// Remover jogador

// Sair da Guild

// Promover Player
