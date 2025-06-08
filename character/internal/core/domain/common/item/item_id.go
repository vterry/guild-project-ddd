package item

import (
	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/base"
)

type ItemID struct {
	base.BaseID[uuid.UUID]
}

func NewItemID(value uuid.UUID) ItemID {
	return ItemID{
		BaseID: base.New(value),
	}
}

func (i ItemID) Equals(otherID ItemID) bool {
	return i.BaseID.Equals(otherID.BaseID)
}
