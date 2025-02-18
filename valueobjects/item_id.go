package valueobjects

import "github.com/google/uuid"

type ItemID struct {
	BaseID[uuid.UUID]
}

func NewItemID(value uuid.UUID) ItemID {
	return ItemID{
		BaseID: NewBaseID(value),
	}
}

func (pID ItemID) Equals(otherID ItemID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
