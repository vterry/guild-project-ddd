package specs

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/common"
	"github.com/vterry/guild-project-ddd/domain/player"
)

const (
	MAX_PLAYERS = 50
	MAX_INVITES = 50
)

var (
	ErrPlayerIsAlreadyGuildMember = errors.New("player is already a guild's member")
	ErrInvalidPlayer              = errors.New("nil referente for player")
	ErrAnotherGuildMember         = errors.New("player is already member of another guild")
	ErrCannotInvite               = errors.New("only guild's member can invite")
	ErrNoInviteAvailable          = errors.New("there are no room for new invitations")
	ErrGuildAlreadyFull           = errors.New("guild is already full")
)

type NewMemberParams struct {
	inviteSender     *player.Player
	guestPlayer      *player.Player
	guildMembers     map[uuid.UUID]*player.Player
	guildInvitesSize int
}

func ValidateNewMember(
	inviteSender *player.Player,
	guestPlayer *player.Player,
	guildMembers map[uuid.UUID]*player.Player,
	guildInvitesSize int,
) error {
	params := NewMemberParams{
		inviteSender:     inviteSender,
		guestPlayer:      guestPlayer,
		guildMembers:     guildMembers,
		guildInvitesSize: guildInvitesSize,
	}
	spec := NewGuildMemberSpecification()
	return spec(common.Base[NewMemberParams]{Entity: &params})
}

// Still reflective if for that example this implementation worth -- It feel It could be implemented in a very simple way

func NewGuildMemberSpecification() common.Specification[NewMemberParams] {
	return common.And(
		SenderMustBePartOfGuild(),
		NotBeingAlreadyMember(),
		GuestNotInAnotherGuildSpec(),
		NoRommForNewSpaceInGuild(),
		NoLeftSpaceInGuild(),
	)
}

func PlayerNotInAnotherGuildSpec[T any](
	getPlayer func(*T) *player.Player,
) common.Specification[T] {

	return func(b common.Base[T]) error {
		if b.Entity == nil {
			return ErrAnotherGuildMember
		}

		player := getPlayer(b.Entity)
		if player == nil {
			return ErrInvalidPlayer
		}

		if guildID := player.GetCurrentGuild(); guildID != "" {
			return ErrAnotherGuildMember
		}

		return nil
	}
}

func SenderMustBePartOfGuild() common.Specification[NewMemberParams] {
	return func(b common.Base[NewMemberParams]) error {
		if _, isMember := b.Entity.guildMembers[b.Entity.inviteSender.ID()]; !isMember {
			return ErrCannotInvite
		}
		return nil
	}
}

func NotBeingAlreadyMember() common.Specification[NewMemberParams] {
	return func(b common.Base[NewMemberParams]) error {
		if _, isMember := b.Entity.guildMembers[b.Entity.guestPlayer.ID()]; isMember {
			return ErrPlayerIsAlreadyGuildMember
		}
		return nil
	}
}

func GuestNotInAnotherGuildSpec() common.Specification[NewMemberParams] {
	return PlayerNotInAnotherGuildSpec(
		func(p *NewMemberParams) *player.Player { return p.guestPlayer },
	)
}

func NoRommForNewSpaceInGuild() common.Specification[NewMemberParams] {
	return func(b common.Base[NewMemberParams]) error {
		if b.Entity.guildInvitesSize+1 >= MAX_INVITES {
			return ErrNoInviteAvailable
		}
		return nil
	}
}

func NoLeftSpaceInGuild() common.Specification[NewMemberParams] {
	return func(b common.Base[NewMemberParams]) error {
		if len(b.Entity.guildMembers)+1 >= MAX_PLAYERS {
			return ErrGuildAlreadyFull
		}
		return nil
	}
}
