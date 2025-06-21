package character

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/base"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/class"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/guild"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/login"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/vault"
	"github.com/vterry/ddd-study/character/internal/core/domain/inventory"
	"github.com/vterry/ddd-study/character/internal/core/domain/playeritem"
)

var (
	ErrCreatePlayer       = errors.New("error while creating player")
	ErrCannotPickItem     = errors.New("error while picking item")
	ErrCannotDropItem     = errors.New("error while dropping item")
	ErrCannotWithdrawGold = errors.New("error while withdrawing gold")
	ErrCannotJoinGuild    = errors.New("character is already member of a guild")
)

type CharacterID struct {
	base.BaseID[uuid.UUID]
}

type Character struct {
	CharacterID
	loginID   login.LoginID
	nickname  string
	class     class.Class
	inventory inventory.Inventory
	guild     guild.GuildID
	vault     vault.VaultID
}

func NewCharacterID(value uuid.UUID) CharacterID {
	return CharacterID{
		BaseID: base.New(value),
	}
}

func CreateNewCharacter(nickname string, loginId login.LoginID, class class.Class, vaultId vault.VaultID) (*Character, error) {

	if err := ValidateNewCharacter(nickname, loginId, class); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreatePlayer, err)
	}

	player := &Character{
		CharacterID: NewCharacterID(uuid.New()),
		loginID:     loginId,
		nickname:    nickname,
		class:       class,
		inventory:   *inventory.NewInventory(),
		guild:       guild.NewGuildID(uuid.Nil),
		vault:       vaultId,
	}

	return player, nil
}

func (c *Character) LoginID() login.LoginID {
	return c.loginID
}

func (c *Character) Nickname() string {
	return c.nickname
}

func (c *Character) Class() class.Class {
	return c.class
}

func (c *Character) Inventory() inventory.Inventory {
	return c.inventory
}

func (c *Character) GetCurrentGuild() guild.GuildID {
	return c.guild
}

func (c *Character) GetCurrentVaultId() vault.VaultID {
	return c.vault
}

func (c *Character) UpdateGuildInfo(guildId guild.GuildID) {
	c.guild = guildId
}

func (c *Character) PickItem(playeritem playeritem.PlayerItem) error {
	if err := c.inventory.AddItem(playeritem); err != nil {
		return fmt.Errorf("%w: %v", ErrCannotPickItem, err)
	}
	return nil
}

func (c *Character) DropItem(playeritem playeritem.PlayerItem) error {
	if err := c.inventory.DropItem(playeritem); err != nil {
		return fmt.Errorf("%w: %v", ErrCannotDropItem, err)
	}
	return nil
}

func (c *Character) PickGold(amount int) error {
	return c.inventory.AddGold(amount)
}

func (c *Character) DropGold(amount int) error {
	if err := c.inventory.WithdrawGold(amount); err != nil {
		return fmt.Errorf("%w: %v", ErrCannotWithdrawGold, err)
	}
	return nil
}

func (c *Character) OpenInventory() []playeritem.PlayerItemID {
	return c.inventory.ShowItems()
}

func (id CharacterID) Equals(other CharacterID) bool {
	return id.BaseID.Equals(other.BaseID)
}
