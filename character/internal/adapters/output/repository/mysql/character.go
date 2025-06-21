package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vterry/ddd-study/character/internal/adapters/output/repository/dao"
	"github.com/vterry/ddd-study/character/internal/core/domain/character"
)

type CharacterRepository struct {
	db *sql.DB
}

func NewCharacterRepository(db *sql.DB) *CharacterRepository {
	return &CharacterRepository{
		db: db,
	}
}

func (c *CharacterRepository) Save(ctx context.Context, character character.Character) error {

	daoInventory := dao.InventorytoDAO(character.Inventory())
	daoCharacter := dao.CharacterToDAO(character)

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting insert transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, CreateNewInventoryQuery, daoInventory.InventoryID, daoInventory.GoldAmount)

	if err != nil {
		return fmt.Errorf("error saving inventory: %w", err)
	}

	_, err = tx.ExecContext(ctx, CreateNewCharacterQuery, daoCharacter.CharacterID, daoCharacter.LoginID, daoCharacter.Nickname, daoCharacter.Class, daoCharacter.InventoryID, daoCharacter.GuildID, daoCharacter.VaultID)

	if err != nil {
		return fmt.Errorf("error saving character: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing saving transaction: %w", err)
	}

	return nil
}

func (c *CharacterRepository) FindCharacterById(ctx context.Context, characterId character.CharacterID) (*character.Character, error) {
	return nil, nil
}

func (c *CharacterRepository) Update(ctx context.Context, character character.Character) error {
	return nil
}
