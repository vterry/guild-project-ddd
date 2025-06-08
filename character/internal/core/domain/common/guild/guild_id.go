package guild

import (
	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/base"
)

type GuildID struct {
	base.BaseID[uuid.UUID]
}

func NewGuildID(value uuid.UUID) GuildID {
	return GuildID{
		BaseID: base.New(value),
	}
}

func (g *GuildID) Equals(otherID GuildID) bool {
	return g.BaseID.Equals(otherID.BaseID)
}
