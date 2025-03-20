package player

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vterry/guild-project-ddd/domain/item"
	"github.com/vterry/guild-project-ddd/domain/player/valueobjects"
)

var (
	defaultName    = "Default Player"
	defaultClass   = valueobjects.Ranger
	playerInstance *Player
	playerOnce     sync.Once
)

func TestInvalidInitialization(t *testing.T) {
	t.Run("test new player with invalid name", func(t *testing.T) {
		_, err := NewPlayer("", valueobjects.Mage)
		assert.ErrorIs(t, err, ErrInvalidNickname)
	})
}

func TestNewPlayer(t *testing.T) {
	t.Cleanup(resetPlayer)
	player := getPlayer()

	t.Run("Test player initialization", func(t *testing.T) {
		assert.NotNil(t, player)
	})

	t.Run("Test player ID is not nil", func(t *testing.T) {
		assert.NotNil(t, player.ID())
	})

	t.Run("Test player default gold value", func(t *testing.T) {
		assert.Equal(t, 0, player.GetCurrentGold())
	})

	t.Run("Test player default cash value", func(t *testing.T) {
		assert.Equal(t, 0, player.GetCurrentCash())
	})

	t.Run("Test player default inventory", func(t *testing.T) {
		assert.Equal(t, 0, len(player.items))
	})

	t.Run("Test player default guild is nil <not a guild member>", func(t *testing.T) {
		assert.Equal(t, "", player.GetCurrentGuild())
	})
}

func TestPlayerAsGuildMember(t *testing.T) {
	t.Cleanup(resetPlayer)
	player := getPlayer()
	mockGuildID := "sampleguild-2006-01-02-15:04:05-129389dsj"

	t.Run("check if player is a guild member", func(t *testing.T) {
		player.UpdateCurrentGuild(mockGuildID)
		assert.NotNil(t, player.currentGuild)
		assert.Equal(t, mockGuildID, player.currentGuild)
	})
}

func TestCashFunctions(t *testing.T) {
	t.Cleanup(resetPlayer)
	player := getPlayer()

	t.Run("Test update cash up to negative value", func(t *testing.T) {
		err := player.UpdateCash(-1)
		assert.ErrorIs(t, err, ErrNotEnoughCash)
		assert.Equal(t, 0, player.GetCurrentCash())
	})

	t.Run("Test update cash up to positive value", func(t *testing.T) {
		err := player.UpdateCash(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, player.GetCurrentCash())
	})
}

func TestGoldFunctions(t *testing.T) {
	t.Cleanup(resetPlayer)
	player := getPlayer()

	t.Run("Test update gold up to negative value", func(t *testing.T) {
		err := player.UpdateGold(-1)
		assert.ErrorIs(t, err, ErrNotEnoughGold)
		assert.Equal(t, 0, player.GetCurrentGold())
	})

	t.Run("Test update gold up to positive value", func(t *testing.T) {
		err := player.UpdateGold(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, player.GetCurrentGold())
	})
}

func TestIventoryFunctions(t *testing.T) {
	playerItem := item.PickRandomItem()

	t.Run("Test retrieve items from empty inventory", func(t *testing.T) {
		t.Cleanup(resetPlayer)
		player := getPlayer()

		err := player.RetriveItem(item.PickRandomItem())
		assert.ErrorIs(t, err, ErrEmptyInventory)
	})

	t.Run("Test adding items to player's inventory", func(t *testing.T) {
		t.Cleanup(resetPlayer)
		player := getPlayer()

		err := player.PickItem(playerItem)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(player.items))
	})

	t.Run("Test adding items to player's inventory with no space left", func(t *testing.T) {
		t.Cleanup(resetPlayer)
		player := getPlayer()
		fullfilPlayerInventory(player)

		err := player.PickItem(playerItem)
		assert.ErrorIs(t, err, ErrNotEnoughSpace)

	})
	t.Run("Test retrieving unexisting item from player's inventory", func(t *testing.T) {
		t.Cleanup(resetPlayer)
		player := getPlayer()
		player.PickItem(playerItem)

		err := player.RetriveItem(item.PickRandomItem())
		assert.ErrorIs(t, err, ErrItemNotFound)
	})

	t.Run("Test retrieving existing item from player's inventory", func(t *testing.T) {
		t.Cleanup(resetPlayer)
		player := getPlayer()
		player.PickItem(playerItem)

		err := player.RetriveItem(playerItem)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(player.items))
	})
}

func fullfilPlayerInventory(player *Player) {
	for i := 0; i <= cap(player.items); i++ {
		player.PickItem(item.PickRandomItem())
	}
}

func resetPlayer() {
	playerInstance, _ = NewPlayer(defaultName, defaultClass)
}

func getPlayer() *Player {
	playerOnce.Do(func() {
		playerInstance, _ = NewPlayer(defaultName, defaultClass)
	})
	return playerInstance
}
