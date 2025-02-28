package valueobjects

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
)

type PlayerID struct {
	common.BaseID[uuid.UUID]
}

func NewPlayerID(value uuid.UUID) PlayerID {
	return PlayerID{
		BaseID: common.NewBaseID(value),
	}
}

func (pID PlayerID) Equals(otherID PlayerID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
