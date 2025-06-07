package vault

import (
	"github.com/google/uuid"
	baseid "github.com/vterry/ddd-study/character/internal/domain/common/base"
)

type VaultID struct {
	baseid.BaseID[uuid.UUID]
}

func New(value uuid.UUID) VaultID {
	return VaultID{
		BaseID: baseid.New(value),
	}
}

func (pID VaultID) Equals(otherID VaultID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
