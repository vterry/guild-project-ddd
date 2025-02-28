package treasure

import (
	"time"

	"github.com/vterry/guild-project-ddd/domain/common"
	player "github.com/vterry/guild-project-ddd/domain/player/valueobjects"
	"github.com/vterry/guild-project-ddd/domain/treasure/valueobjects"
)

type DonationStatus common.Status

type Donation struct {
	valueobjects.DonationID
	playerID  player.PlayerID
	treasure  valueobjects.TreasureID
	status    DonationStatus
	createdAt time.Time
	amount    int
}

func NewDonation(playerId player.PlayerID, treasureId valueobjects.TreasureID, amount int) Donation {
	return Donation{
		playerID:  playerId,
		treasure:  treasureId,
		status:    DonationStatus(common.Pending),
		createdAt: time.Now(),
		amount:    amount,
	}
}
