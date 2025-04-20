package player

import (
	"github.com/google/uuid"
	"github.com/vterry/ddd-study/player/core/common/valueobjects"
)

type PlayerID struct {
	valueobjects.BaseID[uuid.UUID]
}

func NewPlayerID(value uuid.UUID) PlayerID {
	return PlayerID{
		BaseID: valueobjects.NewBaseID(value),
	}
}

func (pID PlayerID) Equals(otherID PlayerID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
