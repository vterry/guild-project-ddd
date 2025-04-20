package valueobjects

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
)

type ItemID struct {
	common.BaseID[uuid.UUID]
}

func NewItemID(value uuid.UUID) ItemID {
	return ItemID{
		BaseID: common.NewBaseID(value),
	}
}

func (pID ItemID) Equals(otherID ItemID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
