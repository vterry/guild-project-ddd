package service

import (
	"github.com/vterry/ddd-study/character/internal/core/domain/character"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/class"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/login"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/vault"
	"github.com/vterry/ddd-study/character/internal/core/domain/playeritem"
)

type CharacterService interface {
	CreateCharacter(loginId login.LoginID, nickname string, class class.Class) error
	TransferItemTo(characterId character.CharacterID, playeritem playeritem.PlayerItem, quantity int, vaultId vault.VaultID) error
	TradeItem(origin character.CharacterID, playeritem playeritem.PlayerItem, quantity int, destiny character.CharacterID) error
	DepositGold(characterId character.CharacterID, quantity int, vaultId vault.VaultID) error
	PickItem() error
	DropItem(playerItemID playeritem.PlayerItemID, quantity int) error
	LeaveGuild(characterID character.CharacterID) error
}
