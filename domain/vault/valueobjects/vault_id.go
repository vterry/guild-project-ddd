package valueobjects

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
)

type VaultID struct {
	common.BaseID[uuid.UUID]
}

func NewVaultID(value uuid.UUID) VaultID {
	return VaultID{
		BaseID: common.NewBaseID(value),
	}
}

func (pID VaultID) Equals(otherID VaultID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
