package guild

import (
	"github.com/google/uuid"
	baseid "github.com/vterry/ddd-study/character/internal/domain/common/base"
)

type GuildID struct {
	baseid.BaseID[uuid.UUID]
}

func New(value uuid.UUID) GuildID {
	return GuildID{
		BaseID: baseid.New(value),
	}
}

func (pID GuildID) Equals(otherID GuildID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
