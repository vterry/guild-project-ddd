package player

import (
	"errors"

	"github.com/vterry/guild-project-ddd/domain/common"
)

var (
	ErrInvalidNickname = errors.New("nickname must by between 4 and 15 characters")
	ErrNotEnoughSpace  = errors.New("player has no space for more items")
	ErrItemNotFound    = errors.New("item not found in player inventory")
	ErrEmptyInventory  = errors.New("player has no items")
	ErrNotEnoughGold   = errors.New("player has not enough gold")
	ErrNotEnoughCash   = errors.New("player has not enough cash")
)

type PlayerError struct {
	common.BaseError
}

func NewPlayerError(curErr, srcErr error) *PlayerError {
	return &PlayerError{
		BaseError: common.NewError(curErr, srcErr).(common.BaseError),
	}
}
