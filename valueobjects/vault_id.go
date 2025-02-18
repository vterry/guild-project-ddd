package valueobjects

import "github.com/google/uuid"

type VaultID struct {
	BaseID[uuid.UUID]
}

func NewVaultID(value uuid.UUID) VaultID {
	return VaultID{
		BaseID: NewBaseID(value),
	}
}

func (pID VaultID) Equals(otherID VaultID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
