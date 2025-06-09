package character

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/class"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/guild"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/item"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/login"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/vault"
	"github.com/vterry/ddd-study/character/internal/core/domain/playeritem"
)

// setupTestCharacter creates a new character for testing
func setupTestCharacter(t *testing.T) *Character {
	validLogin := login.NewLoginID(uuid.New())

	character, err := CreateNewCharacter("TestPlayer", validLogin, class.Warrior, vault.NewVaultID(uuid.New()))
	assert.NoError(t, err)
	assert.NotNil(t, character)

	return character
}

// setupTestItem creates a new item for testing
func setupTestItem(t *testing.T, name string) *playeritem.PlayerItem {
	itemId := item.NewItemID(uuid.New())
	testItem, err := playeritem.NewPlayerItem(itemId, name, 1)
	assert.NoError(t, err)
	assert.NotNil(t, testItem)

	return testItem
}

func TestCreateNewCharacter(t *testing.T) {
	validLogin := login.NewLoginID(uuid.New())

	tests := []struct {
		name      string
		nickname  string
		loginId   login.LoginID
		class     class.Class
		vaultId   vault.VaultID
		wantErr   bool
		errString string
	}{
		{
			name:     "successful character creation",
			nickname: "TestPlayer",
			loginId:  validLogin,
			class:    class.Warrior,
			vaultId:  vault.NewVaultID(uuid.New()),
			wantErr:  false,
		},
		{
			name:      "empty nickname",
			nickname:  "",
			loginId:   validLogin,
			class:     class.Warrior,
			vaultId:   vault.NewVaultID(uuid.New()),
			wantErr:   true,
			errString: "error while creating player",
		},
		{
			name:      "nil login",
			nickname:  "TestPlayer",
			loginId:   login.LoginID{},
			class:     class.Warrior,
			vaultId:   vault.NewVaultID(uuid.New()),
			wantErr:   true,
			errString: "error while creating player",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			character, err := CreateNewCharacter(tt.nickname, tt.loginId, tt.class, tt.vaultId)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errString)
				assert.Nil(t, character)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, character)
				assert.Equal(t, tt.nickname, character.nickname)
				assert.Equal(t, tt.class, character.class)
				assert.Equal(t, tt.vaultId, character.vault)
				assert.NotNil(t, character.inventory)
				assert.Equal(t, guild.NewGuildID(uuid.Nil), character.guild)
			}
		})
	}
}

func TestCharacterGuildOperations(t *testing.T) {
	character := setupTestCharacter(t)

	t.Run("get current guild", func(t *testing.T) {
		guildId := character.GetCurrentGuild()
		assert.Equal(t, guild.NewGuildID(uuid.Nil), guildId)
	})

	t.Run("update guild info", func(t *testing.T) {
		newGuildId := guild.NewGuildID(uuid.New())
		character.UpdateGuildInfo(newGuildId)
		assert.Equal(t, newGuildId, character.GetCurrentGuild())
	})
}

func TestCharacterInventoryOperations(t *testing.T) {
	t.Run("pick and drop item successfully", func(t *testing.T) {
		character := setupTestCharacter(t)
		testItem := setupTestItem(t, "TestItem")

		// Test picking item
		err := character.PickItem(*testItem)
		assert.NoError(t, err)

		// Verify item is in inventory
		items := character.OpenInventory()
		assert.Len(t, items, 1)

		// Test dropping item
		err = character.DropItem(*testItem)
		assert.NoError(t, err)

		// Verify inventory is empty
		items = character.OpenInventory()
		assert.Len(t, items, 0)
	})

	t.Run("pick item with full inventory", func(t *testing.T) {
		character := setupTestCharacter(t)

		// Fill inventory to max capacity
		for i := 0; i < 10; i++ {
			testItem := setupTestItem(t, fmt.Sprintf("TestItem%d", i))
			err := character.PickItem(*testItem)
			assert.NoError(t, err)
		}

		// Verify inventory is full
		items := character.OpenInventory()
		assert.Len(t, items, 10)

		// Try to add one more item
		extraItem := setupTestItem(t, "ExtraItem")
		err := character.PickItem(*extraItem)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot add item - inventory is full")

		// Verify inventory still has only 10 items
		items = character.OpenInventory()
		assert.Len(t, items, 10)
	})

	t.Run("drop non-existent item", func(t *testing.T) {
		character := setupTestCharacter(t)
		nonExistentItem := setupTestItem(t, "NonExistentItem")

		// Try to drop the non-existent item
		err := character.DropItem(*nonExistentItem)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot drop item - player item is not in inventory")

		// Verify inventory is still empty
		items := character.OpenInventory()
		assert.Len(t, items, 0)
	})

	t.Run("gold operations", func(t *testing.T) {
		character := setupTestCharacter(t)

		// Test adding gold
		err := character.PickGold(100)
		assert.NoError(t, err)

		// Test withdrawing gold
		err = character.DropGold(50)
		assert.NoError(t, err)

		// Test withdrawing more than available
		err = character.DropGold(100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error while withdrawing gold")
	})
}

func TestCharacterVaultOperations(t *testing.T) {
	vaultId := vault.NewVaultID(uuid.New())
	validLogin := login.NewLoginID(uuid.New())
	character, _ := CreateNewCharacter("TestPlayer", validLogin, class.Warrior, vaultId)

	t.Run("get vault id", func(t *testing.T) {
		assert.Equal(t, vaultId, character.GetCurrentVaultId())
	})
}
