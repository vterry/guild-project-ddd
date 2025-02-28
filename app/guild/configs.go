package service

import (
	guild "github.com/vterry/guild-project-ddd/app/guild/repository"
	player "github.com/vterry/guild-project-ddd/app/player/repository"
)

type GuildServiceConfiguration func(gs *GuildService) error

func setupGuildRepository(gr guild.GuildRepository) GuildServiceConfiguration {
	return func(gs *GuildService) error {
		gs.guilds = gr
		return nil
	}
}

func setupPlayerRepository(pl player.PlayerRepository) GuildServiceConfiguration {
	return func(gs *GuildService) error {
		gs.players = pl
		return nil
	}
}
