package valueobjects

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
)

type GuildID struct {
	common.BaseID[uuid.UUID]
}

func NewGuildID(value uuid.UUID) GuildID {
	return GuildID{
		BaseID: common.NewBaseID(value),
	}
}

func (pID GuildID) Equals(otherID GuildID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
