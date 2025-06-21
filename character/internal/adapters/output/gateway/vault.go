package gateway

import (
	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/vault"
)

// TODO - implement

type MockVaultGateway struct {
}

func NewMockVaultGateway() *MockVaultGateway {
	return &MockVaultGateway{}
}

func (v *MockVaultGateway) CreateVault() (vault.VaultID, error) {
	return vault.NewVaultID(uuid.New()), nil
}
