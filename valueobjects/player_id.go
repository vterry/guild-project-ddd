package valueobjects

import "github.com/google/uuid"

type PlayerID struct {
	BaseID[uuid.UUID]
}

func NewPlayerID(value uuid.UUID) PlayerID {
	return PlayerID{
		BaseID: NewBaseID(value),
	}
}

func (pID PlayerID) Equals(otherID PlayerID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
