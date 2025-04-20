package guild

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
)

type DonationID struct {
	common.BaseID[uuid.UUID]
}

func NewDonationID(value uuid.UUID) DonationID {
	return DonationID{
		BaseID: common.NewBaseID(value),
	}
}

func (pID DonationID) Equals(otherID DonationID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
