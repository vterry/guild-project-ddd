package valueobjects

import (
	"fmt"

	"github.com/vterry/guild-project-ddd/domain/common"
)

type TreasureID struct {
	common.BaseID[string]
}

func NewTreasureID(value string) TreasureID {
	rand, _ := common.ShortUUID(8)
	treasureId := fmt.Sprintf("%s-%s", value, rand)
	return TreasureID{
		BaseID: common.NewBaseID(treasureId),
	}
}

func (pID TreasureID) Equals(otherID TreasureID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
