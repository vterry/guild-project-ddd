package vault

import (
	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/base"
)

type VaultID struct {
	base.BaseID[uuid.UUID]
}

func NewVaultID(value uuid.UUID) VaultID {
	return VaultID{
		BaseID: base.New(value),
	}
}

func (v VaultID) Equals(otherID VaultID) bool {
	return v.BaseID.Equals(otherID.BaseID)
}
