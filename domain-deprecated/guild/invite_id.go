package guild

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
)

type InviteID struct {
	common.BaseID[uuid.UUID]
}

func NewInviteID(value uuid.UUID) InviteID {
	return InviteID{
		BaseID: common.NewBaseID(value),
	}
}

func (pID InviteID) Equals(otherID InviteID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
