package valueobjects

import (
	"fmt"
	"time"

	"github.com/vterry/guild-project-ddd/domain/common"
)

type GuildID struct {
	common.BaseID[string]
}

func NewGuildID(guildName string, createdAt time.Time) GuildID {
	rand, _ := common.ShortUUID(8)
	guildId := fmt.Sprintf("%s-%s-%s", guildName, createdAt.Format("2006-01-02-150405"), rand)
	return GuildID{
		BaseID: common.NewBaseID(guildId),
	}
}

func (pID GuildID) Equals(otherID GuildID) bool {
	return pID.BaseID.Equals(otherID.BaseID)
}
