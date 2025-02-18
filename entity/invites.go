package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/valueobjects"
)

type InviteStatus string

const (
	Pending  InviteStatus = "PENDING"
	Accepted InviteStatus = "ACCEPTED"
	Rejected InviteStatus = "REJECTED"
)

type Invite struct {
	valueobjects.InviteID
	playerID  valueobjects.PlayerID
	invitedBy valueobjects.PlayerID
	guildID   valueobjects.GuildID
	status    InviteStatus
	createdAt time.Time
}

func NewInvite(invited valueobjects.PlayerID, sender valueobjects.PlayerID, guild valueobjects.GuildID) Invite {
	return Invite{
		InviteID:  valueobjects.NewInviteID(uuid.New()),
		playerID:  invited,
		invitedBy: sender,
		guildID:   guild,
		status:    Pending,
		createdAt: time.Now(),
	}
}

func (i Invite) GetPlayerID() valueobjects.PlayerID {
	return i.playerID
}

func (i Invite) CheckStatus() InviteStatus {
	return i.status
}

func (i *Invite) UpdateStatus(status InviteStatus) InviteStatus {
	i.status = status
	return i.status
}
