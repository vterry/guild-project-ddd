package playeritem

import (
	"github.com/google/uuid"
	baseid "github.com/vterry/ddd-study/character/internal/domain/common/base"
)

type PlayerItemID struct {
	baseid.BaseID[uuid.UUID]
}

func NewPlayerItemID(value uuid.UUID) PlayerItemID {
	return PlayerItemID{
		BaseID: baseid.New(value),
	}
}

func (pID PlayerItemID) Equals(otherID PlayerItemID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
