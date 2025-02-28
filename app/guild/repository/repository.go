package repository

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/guild"
)

type GuildRepository interface {
	Get(uuid.UUID) (guild.Guild, error)
	Add(guild.Guild) error
	Update(guild.Guild) error
}
