package character

import (
	"github.com/google/uuid"
	baseid "github.com/vterry/ddd-study/character/internal/domain/common/base"
)

type CharacterID struct {
	baseid.BaseID[uuid.UUID]
}

func CreateID(value uuid.UUID) CharacterID {
	return CharacterID{
		BaseID: baseid.New(value),
	}
}

func (pID CharacterID) Equals(otherID CharacterID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
