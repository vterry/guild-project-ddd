package guild

import (
	"time"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
	"github.com/vterry/guild-project-ddd/domain/guild/valueobjects"
	player "github.com/vterry/guild-project-ddd/domain/player/valueobjects"
)

type InviteStatus common.Status

type Invite struct {
	valueobjects.InviteID
	playerID  player.PlayerID
	invitedBy player.PlayerID
	guildID   valueobjects.GuildID
	status    InviteStatus
	createdAt time.Time
}

func NewInvite(invited player.PlayerID, sender player.PlayerID, guild valueobjects.GuildID) *Invite {
	return &Invite{
		InviteID:  valueobjects.NewInviteID(uuid.New()),
		playerID:  invited,
		invitedBy: sender,
		guildID:   guild,
		status:    InviteStatus(common.Pending),
		createdAt: time.Now(),
	}
}

func (i Invite) GetPlayerID() player.PlayerID {
	return i.playerID
}

func (i Invite) CheckStatus() InviteStatus {
	return i.status
}

func (i *Invite) UpdateStatus(status InviteStatus) InviteStatus {
	i.status = status
	return i.status
}

func (i *Invite) reject() error {
	if i.status != InviteStatus(common.Pending) {
		return ErrInvalidOperation
	}
	i.status = InviteStatus(common.Rejected)
	return nil
}

func (i *Invite) cancel() error {
	if i.status != InviteStatus(common.Pending) {
		return ErrInvalidOperation
	}
	i.status = InviteStatus(common.Canceled)
	return nil
}

func (i *Invite) approve() error {
	if i.status != InviteStatus(common.Pending) {
		return ErrInvalidOperation
	}
	i.status = InviteStatus(common.Approved)
	return nil
}
