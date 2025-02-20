package valueobjects

import "github.com/google/uuid"

type DonationID struct {
	BaseID[uuid.UUID]
}

func NewDonationID(value uuid.UUID) DonationID {
	return DonationID{
		BaseID: NewBaseID(value),
	}
}

func (pID DonationID) Equals(otherID DonationID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
