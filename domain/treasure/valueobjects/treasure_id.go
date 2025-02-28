package valueobjects

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
)

type TreasureID struct {
	common.BaseID[uuid.UUID]
}

func NewTreasureID(value uuid.UUID) TreasureID {
	return TreasureID{
		BaseID: common.NewBaseID(value),
	}
}

func (pID TreasureID) Equals(otherID TreasureID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
