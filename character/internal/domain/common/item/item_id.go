package item

import (
	"github.com/google/uuid"
	baseid "github.com/vterry/ddd-study/character/internal/domain/common/base"
)

type ItemID struct {
	baseid.BaseID[uuid.UUID]
}

func New(value uuid.UUID) ItemID {
	return ItemID{
		BaseID: baseid.New(value),
	}
}

func (pID ItemID) Equals(otherID ItemID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
