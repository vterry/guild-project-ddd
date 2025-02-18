package valueobjects

import "github.com/google/uuid"

type TreasureID struct {
	BaseID[uuid.UUID]
}

func NewTreasureID(value uuid.UUID) TreasureID {
	return TreasureID{
		BaseID: NewBaseID(value),
	}
}

func (pID TreasureID) Equals(otherID TreasureID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
