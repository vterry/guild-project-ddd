package character

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/domain/common/class"
	"github.com/vterry/ddd-study/character/internal/domain/common/guild"
	"github.com/vterry/ddd-study/character/internal/domain/common/vault"
	"github.com/vterry/ddd-study/character/internal/domain/inventory"
	"github.com/vterry/ddd-study/character/internal/domain/login"
	"github.com/vterry/ddd-study/character/internal/domain/playeritem"
)

var (
	ErrCreatePlayer       = errors.New("error while creating player")
	ErrCannotPickItem     = errors.New("error while picking item")
	ErrCannotDropItem     = errors.New("error while dropping item")
	ErrCannotWithdrawGold = errors.New("error while withdrawing gold")
	ErrCannotJoinGuild    = errors.New("character is already member of a guild")
)

type Character struct {
	CharacterID
	login     *login.Login
	nickname  string
	class     class.Class
	inventory inventory.Inventory
	guild     guild.GuildID
	vault     vault.VaultID
}

func CreateNewCharacter(nickname string, login *login.Login, class class.Class, vaultId vault.VaultID) (*Character, error) {

	if err := ValidateNewCharacter(nickname, login, class); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreatePlayer, err)
	}

	player := &Character{
		CharacterID: CreateID(uuid.New()),
		login:       login,
		nickname:    nickname,
		class:       class,
		inventory:   *inventory.NewInventory(),
		guild:       guild.New(uuid.Nil),
		vault:       vaultId,
	}

	return player, nil
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

func (c *Character) PickGold(amount int) {
	c.inventory.AddGold(amount)
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
