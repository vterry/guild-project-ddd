package service

import (
	"context"

	"github.com/vterry/ddd-study/character/internal/core/domain/character"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/class"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/login"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/vault"
	"github.com/vterry/ddd-study/character/internal/core/domain/playeritem"
)

type CharacterService interface {
	CreateCharacter(ctx context.Context, loginId login.LoginID, nickname string, class class.Class) error
	TransferItemTo(ctx context.Context, characterId character.CharacterID, playeritem playeritem.PlayerItem, quantity int, vaultId vault.VaultID) error
	TradeItem(ctx context.Context, origin character.CharacterID, playeritem playeritem.PlayerItem, quantity int, destiny character.CharacterID) error
	DepositGold(ctx context.Context, characterId character.CharacterID, quantity int, vaultId vault.VaultID) error
	PickItem() error
	DropItem(ctx context.Context, playerItemID playeritem.PlayerItemID, quantity int) error
	LeaveGuild(ctx context.Context, characterID character.CharacterID) error
}
