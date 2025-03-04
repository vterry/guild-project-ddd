package guild

import (
	"errors"

	"github.com/vterry/guild-project-ddd/domain/common"
)

var (
	ErrInvalidGuildName           = errors.New("guild's name must by between 4 and 15 characters")
	ErrMustInformGuidOwner        = errors.New("a guild master must be inform")
	ErrPlayerIsAlreadyGuildMember = errors.New("player is already a guild's member")
	ErrAnotherGuildMember         = errors.New("player is already member of another guild")
	ErrCannotInvite               = errors.New("only guild's member can invite")
	ErrNotGuildMember             = errors.New("player isnt a guild member")
	ErrAlreadyInvited             = errors.New("this player has already a pending invite")
	ErrInviteNotExistis           = errors.New("invite no longer exists")
	ErrInvalidInviteState         = errors.New("invite has an invalid state")
	ErrInvalidOperation           = errors.New("invalid operation")
	ErrNoInviteAvailable          = errors.New("there are no room for new invitations")
	ErrGuildAlreadyFull           = errors.New("guild is already full")
)

type GuildError struct {
	common.BaseError
}

func NewGuildError(curErr, srcErr error) *GuildError {
	return &GuildError{
		BaseError: common.NewError(curErr, srcErr).(common.BaseError),
	}
}
