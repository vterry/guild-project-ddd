package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/valueobjects"
)

type InviteStatus valueobjects.Status

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
		status:    InviteStatus(valueobjects.Pending),
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
