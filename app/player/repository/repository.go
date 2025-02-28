package repository

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/player"
)

type PlayerRepository interface {
	Get(uuid.UUID) (player.Player, error)
	Add(player.Player) error
	Update(player.Player) error
}
