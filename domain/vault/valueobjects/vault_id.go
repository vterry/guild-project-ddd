package valueobjects

import (
	"fmt"

	"github.com/vterry/guild-project-ddd/domain/common"
)

type VaultID struct {
	common.BaseID[string]
}

func NewVaultID(value string) VaultID {
	rand, _ := common.ShortUUID(8)
	vaultId := fmt.Sprintf("%s-%s", value, rand)
	return VaultID{
		BaseID: common.NewBaseID(vaultId),
	}
}

func (pID VaultID) Equals(otherID VaultID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
