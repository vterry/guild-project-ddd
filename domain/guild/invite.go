package guild

import (
	"time"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
	"github.com/vterry/guild-project-ddd/domain/guild/valueobjects"
)

type InviteStatus common.Status

type Invite struct {
	valueobjects.InviteID
	playerID  uuid.UUID
	invitedBy uuid.UUID
	guildID   uuid.UUID
	status    InviteStatus
	createdAt time.Time
}

func NewInvite(guest uuid.UUID, sender uuid.UUID, guild uuid.UUID) *Invite {
	return &Invite{
		InviteID:  valueobjects.NewInviteID(uuid.New()),
		playerID:  guest,
		invitedBy: sender,
		guildID:   guild,
		status:    InviteStatus(common.Pending),
		createdAt: time.Now(),
	}
}

func (i Invite) GetPlayerID() uuid.UUID {
	return i.playerID
}

func (i Invite) CheckStatus() InviteStatus {
	return i.status
}

func (i *Invite) reject() error {
	if i.status != InviteStatus(common.Pending) {
		return ErrInvalidInviteState
	}
	i.status = InviteStatus(common.Rejected)
	return nil
}

func (i *Invite) cancel() error {
	if i.status != InviteStatus(common.Pending) {
		return ErrInvalidInviteState
	}
	i.status = InviteStatus(common.Canceled)
	return nil
}

func (i *Invite) approve() error {
	if i.status != InviteStatus(common.Pending) {
		return ErrInvalidInviteState
	}
	i.status = InviteStatus(common.Approved)
	return nil
}
