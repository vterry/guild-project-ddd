package guild

import (
	"time"

	"github.com/vterry/guild-project-ddd/domain/common"
	"github.com/vterry/guild-project-ddd/domain/player"
)

type DonationStatus common.Status

type Donation struct {
	DonationID
	playerID  player.PlayerID
	status    DonationStatus
	createdAt time.Time
	amount    int
}

func NewDonation(playerId player.PlayerID, amount int) Donation {
	return Donation{
		playerID:  playerId,
		status:    DonationStatus(common.Pending),
		createdAt: time.Now(),
		amount:    amount,
	}
}
