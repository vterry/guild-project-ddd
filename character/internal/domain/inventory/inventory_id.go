package inventory

import (
	"github.com/google/uuid"
	baseid "github.com/vterry/ddd-study/character/internal/domain/common/base"
)

type InventoryID struct {
	baseid.BaseID[uuid.UUID]
}

func New(value uuid.UUID) InventoryID {
	return InventoryID{
		BaseID: baseid.New(value),
	}
}

func (pID InventoryID) Equals(otherID InventoryID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
