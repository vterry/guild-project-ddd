package valueobjects

import "github.com/google/uuid"

type GuildID struct {
	BaseID[uuid.UUID]
}

func NewGuildID(value uuid.UUID) GuildID {
	return GuildID{
		BaseID: NewBaseID(value),
	}
}

func (pID GuildID) Equals(otherID GuildID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
