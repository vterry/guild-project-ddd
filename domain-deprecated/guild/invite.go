package guild

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
	"github.com/vterry/guild-project-ddd/domain/common/valueobjects"
)

var (
	ErrInvalidInviteState = errors.New("invite has an invalid state")
)

type InviteStatus common.Status

type Invite struct {
	valueobjects.InviteID
	playerID  uuid.UUID
	invitedBy uuid.UUID
	guildID   string
	status    InviteStatus
	createdAt time.Time
}

func NewInvite(guestId uuid.UUID, senderId uuid.UUID, guildId string) *Invite {
	return &Invite{
		InviteID:  valueobjects.NewInviteID(uuid.New()),
		playerID:  guestId,
		invitedBy: senderId,
		guildID:   guildId,
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
