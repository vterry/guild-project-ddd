package vault

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vterry/guild-project-ddd/domain/common"
	"github.com/vterry/guild-project-ddd/domain/item"
	"github.com/vterry/guild-project-ddd/domain/player"
)

var (
	vaultInstance *Vault
	vaultOnce     sync.Once
)

func TestInicializeVault(t *testing.T) {
	t.Cleanup(resetVault)
	vault := getVault()

	if vault == nil {
		t.Error("Error initializing vault")
		return
	}

	t.Run("Check initial vault gold amount", func(t *testing.T) {
		assert.Equal(t, 0, vault.GoldAmount)
	})

	t.Run("Check initial items stored", func(t *testing.T) {
		assert.Equal(t, 0, len(vault.Items))
	})

}

func TestVaultID(t *testing.T) {
	vault := NewVault()
	if vault == nil {
		t.Error("Error initializing vault")
		return
	}

	t.Run("Check vault ID", func(t *testing.T) {
		assert.NotEqual(t, "", vault.ID())
	})
}

func TestAddItem(t *testing.T) {
	t.Cleanup(resetVault)
	vault := NewVault()

	player, _ := player.NewPlayer("Player 1", common.Warrior)
	item1 := item.PickRandomItem()
	item2 := item.PickRandomItem()
	player.PickItem(item1)

	if vault == nil {
		t.Error("Error initializing vault")
		return
	}

	t.Run("Add item that player not have", func(t *testing.T) {
		var vaultError *VaultError
		err := vault.AddItem(item2, player)
		assert.ErrorAs(t, err, &vaultError)
		assert.Contains(t, ErrInvalidOperation.Error(), strings.Split(vaultError.Error(), "\n")[0])
		assert.Equal(t, 0, len(vault.Items))
	})

	t.Run("Add item that player has", func(t *testing.T) {
		vault.AddItem(item1, player)
		assert.NotEqual(t, 0, len(vault.Items))
		assert.Positive(t, len(vault.Items))
		assert.Equal(t, item1.Name(), vault.Items[0].Name())
	})
}

func TestRetriveItem(t *testing.T) {
	t.Cleanup(resetVault)
	vault := NewVault()

	regPlayer, _ := player.NewPlayer("RegularPlayer", common.Ranger)
	maxInvPlayer, _ := player.NewPlayer("MaxInvPlayer", common.Warrior)

	item1 := item.PickRandomItem()
	item2 := item.PickRandomItem()
	item3 := item.PickRandomItem()

	fillPlayerInventory(regPlayer, item1, 2)
	fillPlayerInventory(maxInvPlayer, item2, 10)

	if vault == nil {
		t.Error("Error initializing vault")
		return
	}

	t.Run("Retrieve item from empty vault", func(t *testing.T) {
		var vaultError *VaultError
		err := vault.RetriveItem(item1, regPlayer)
		assert.ErrorAs(t, err, &vaultError)
		assert.Contains(t, ErrEmptyVault.Error(), vaultError.Error())
	})

	t.Run("Retrieve existing item from vault", func(t *testing.T) {
		var vaultError *VaultError

		vault.AddItem(item1, regPlayer)
		assert.Equal(t, 1, len(vault.Items))

		err := vault.RetriveItem(item1, regPlayer)
		assert.NotErrorAs(t, err, &vaultError)
		assert.Equal(t, 0, len(vault.Items))
	})

	t.Run("Retrieve unexisting item from vault", func(t *testing.T) {
		var vaultErr *VaultError

		vault.AddItem(item1, regPlayer)
		vaultSize := len(vault.Items)

		err := vault.RetriveItem(item.PickRandomItem(), regPlayer)
		assert.Error(t, vault.RetriveItem(item3, regPlayer))

		assert.ErrorAs(t, err, &vaultErr)
		assert.Contains(t, vaultErr.Error(), ErrItemNotFound.Error())
		assert.Equal(t, vaultSize, len(vault.Items))

	})

	t.Run("Player retriving item from vault with no space left in inventory", func(t *testing.T) {
		var vaultError *VaultError

		vault.AddItem(item1, regPlayer)
		vaultSize := len(vault.Items)
		err := vault.RetriveItem(item1, maxInvPlayer)

		assert.ErrorAs(t, err, &vaultError)
		assert.Equal(t, vaultSize, len(vault.Items))
	})

}

func TestAddGold(t *testing.T) {
	t.Cleanup(resetVault)
	vault := NewVault()
	player, _ := player.NewPlayer("Player 1", common.Ranger)

	t.Run("Deposit negative gold amount", func(t *testing.T) {
		var vaultError *VaultError

		err := vault.AddGold(-1, player)
		assert.ErrorAs(t, err, &vaultError)
		assert.Equal(t, ErrNegativeGoldAmount.Error(), vaultError.Error())
	})

	t.Run("Deposit valid amount of gold", func(t *testing.T) {
		player.UpdateGold(1000)
		err := vault.AddGold(500, player)
		assert.NoError(t, err)
		assert.Equal(t, vault.GoldAmount, 500)
	})

	t.Run("Deposit value that player doestn have", func(t *testing.T) {
		var vaultError *VaultError

		err := vault.AddGold(2000, player)
		assert.ErrorAs(t, err, &vaultError)
		assert.Equal(t, ErrInvalidGoldAmount.Error(), vaultError.Error())
	})

}

func TestGoldWithdraw(t *testing.T) {
	t.Cleanup(resetVault)
	vault := NewVault()
	player, _ := player.NewPlayer("Player 1", common.Ranger)

	t.Run("Withdraw negative value", func(t *testing.T) {
		var vaultError *VaultError

		err := vault.GoldWithdraw(-1, player)
		assert.ErrorAs(t, err, &vaultError)
		assert.Equal(t, ErrNegativeGoldAmount.Error(), vaultError.Error())
	})

	t.Run("Withdraw negative that vault dont have", func(t *testing.T) {
		var vaultError *VaultError

		err := vault.GoldWithdraw(1, player)
		assert.ErrorAs(t, err, &vaultError)
		assert.Equal(t, ErrInvalidGoldAmount.Error(), vaultError.Error())
	})

	t.Run("Withdraw correct value from vault", func(t *testing.T) {
		player.UpdateGold(1000)
		vault.AddGold(500, player)
		err := vault.GoldWithdraw(500, player)
		assert.NoError(t, err)
		assert.Equal(t, vault.GoldAmount, 0)
		assert.Equal(t, player.GetCurrentGold(), 1000)
	})
}

func fillPlayerInventory(p *player.Player, i *item.Item, t int) {
	for j := 0; j <= t; j++ {
		p.PickItem(i)
	}
}

func resetVault() {
	vaultInstance = NewVault()
}

func getVault() *Vault {
	vaultOnce.Do(func() {
		vaultInstance = NewVault()
	})
	return vaultInstance
}
