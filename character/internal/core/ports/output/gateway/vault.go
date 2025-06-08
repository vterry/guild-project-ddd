package gateway

import "github.com/vterry/ddd-study/character/internal/core/domain/common/vault"

type Vault interface {
	CreateVault() (vault.VaultID, error)
}
