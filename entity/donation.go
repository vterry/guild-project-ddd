package entity

import (
	"time"

	"github.com/vterry/guild-project-ddd/valueobjects"
)

type DonationStatus valueobjects.Status

type Donation struct {
	valueobjects.DonationID
	playerID  valueobjects.PlayerID
	treasure  valueobjects.TreasureID
	status    DonationStatus
	createdAt time.Time
	amount    int
}

func NewDonation(playerId valueobjects.PlayerID, treasureId valueobjects.TreasureID, amount int) Donation {
	return Donation{
		playerID:  playerId,
		treasure:  treasureId,
		status:    DonationStatus(valueobjects.Pending),
		createdAt: time.Now(),
		amount:    amount,
	}
}
