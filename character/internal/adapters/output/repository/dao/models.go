package dao

import (
	"github.com/vterry/ddd-study/character/internal/core/domain/character"
	"github.com/vterry/ddd-study/character/internal/core/domain/inventory"
	"github.com/vterry/ddd-study/character/internal/core/domain/playeritem"
)

type Character struct {
	CharacterID string
	LoginID     string
	Nickname    string
	Class       string
	InventoryID string //understand how to properly set up the relationship
	GuildID     string
	VaultID     string
}

type Inventory struct {
	InventoryID string
	GoldAmount  int
	Items       map[playeritem.PlayerItemID]playeritem.PlayerItem //understand how to properly set up the relationship - One to Many
}

type PlayerItem struct {
	PlayerItemId string
	ItemID       string
	Description  string
	Quantity     int
}

func InventorytoDAO(inventory inventory.Inventory) *Inventory {
	return &Inventory{
		InventoryID: inventory.ID().String(),
		GoldAmount:  inventory.GetCurrentGold(),
	}
}

func CharacterToDAO(character character.Character) *Character {
	return &Character{
		CharacterID: character.ID().String(),
		LoginID:     character.LoginID().ID().String(),
		Nickname:    character.Nickname(),
		Class:       character.Class().String(),
		InventoryID: character.Inventory().ID().String(),
		GuildID:     character.GetCurrentGuild().ID().String(),
		VaultID:     character.GetCurrentVaultId().ID().String(),
	}
}

func DAOToCharacter(dao *Character) *character.Character {
	return &character.Character{}
}
