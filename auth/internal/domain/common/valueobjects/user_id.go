package valueobjects

import "github.com/google/uuid"

type UserID struct {
	BaseID[uuid.UUID]
}

func NewUserID(value uuid.UUID) UserID {
	return UserID{
		BaseID: NewBaseID(value),
	}
}

func (pID UserID) Equals(otherID UserID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
