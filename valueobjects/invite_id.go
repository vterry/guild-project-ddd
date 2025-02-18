package valueobjects

import "github.com/google/uuid"

type InviteID struct {
	BaseID[uuid.UUID]
}

func NewInviteID(value uuid.UUID) InviteID {
	return InviteID{
		BaseID: NewBaseID(value),
	}
}

func (pID InviteID) Equals(otherID InviteID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
