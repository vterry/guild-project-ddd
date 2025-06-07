package inventory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vterry/ddd-study/character/internal/domain/common/item"
	"github.com/vterry/ddd-study/character/internal/domain/playeritem"
)

func TestNewInventory(t *testing.T) {
	inv := NewInventory()
	assert.NotNil(t, inv)
	assert.Equal(t, 0, inv.goldAmount)
	assert.NotNil(t, inv.items)
	assert.Equal(t, 0, len(inv.items))
}

func TestAddItem(t *testing.T) {
	tests := []struct {
		name          string
		setupItems    int
		shouldSucceed bool
	}{
		{"add item to empty inventory", 0, true},
		{"add item to partially filled inventory", 5, true},
		{"add item to full inventory", MAX_ITEMS, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := NewInventory()

			// Fill inventory up to setupItems
			for i := 0; i < tt.setupItems; i++ {
				itemID := item.New(uuid.New())
				playerItem, err := playeritem.NewPlayerItem(itemID, "Test Item", 1)
				assert.NoError(t, err)
				err = inv.AddItem(*playerItem)
				assert.NoError(t, err)
			}

			// Try to add one more item
			itemID := item.New(uuid.New())
			newItem, err := playeritem.NewPlayerItem(itemID, "Test Item", 1)
			assert.NoError(t, err)
			err = inv.AddItem(*newItem)

			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.Equal(t, tt.setupItems+1, len(inv.items))
				assert.Contains(t, inv.items, newItem.PlayerItemID)
			} else {
				assert.Error(t, err)
				assert.Equal(t, ErrInventoryIsFull, err)
				assert.Equal(t, tt.setupItems, len(inv.items))
			}
		})
	}
}

func TestDropItem(t *testing.T) {
	tests := []struct {
		name          string
		setupItems    int
		shouldSucceed bool
	}{
		{"drop item from inventory with one item", 1, true},
		{"drop item from empty inventory", 0, false},
		{"drop non-existent item", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := NewInventory()
			var testItem *playeritem.PlayerItem

			// Setup inventory
			if tt.setupItems > 0 {
				itemID := item.New(uuid.New())
				var err error
				testItem, err = playeritem.NewPlayerItem(itemID, "Test Item", 1)
				assert.NoError(t, err)
				err = inv.AddItem(*testItem)
				assert.NoError(t, err)
			}

			// Try to drop item
			if tt.name == "drop non-existent item" {
				itemID := item.New(uuid.New())
				nonExistentItem, err := playeritem.NewPlayerItem(itemID, "Test Item", 1)
				assert.NoError(t, err)
				err = inv.DropItem(*nonExistentItem)
				assert.Error(t, err)
				assert.Equal(t, ErrPlayerItemNotFound, err)
			} else if tt.name == "drop item from empty inventory" {
				// Create a new item but don't add it to inventory
				itemID := item.New(uuid.New())
				emptyItem, err := playeritem.NewPlayerItem(itemID, "Test Item", 1)
				assert.NoError(t, err)
				err = inv.DropItem(*emptyItem)
				assert.Error(t, err)
				assert.Equal(t, ErrPlayerItemNotFound, err)
			} else {
				err := inv.DropItem(*testItem)
				if tt.shouldSucceed {
					assert.NoError(t, err)
					assert.Equal(t, 0, len(inv.items))
				} else {
					assert.Error(t, err)
					assert.Equal(t, ErrPlayerItemNotFound, err)
				}
			}
		})
	}
}

func TestAddGold(t *testing.T) {
	tests := []struct {
		name          string
		amount        int
		shouldSucceed bool
		expectedGold  int
	}{
		{"add positive amount", 100, true, 100},
		{"add zero amount", 0, true, 0},
		{"add negative amount", -50, false, 0},
		{"add large amount", 999999, true, 999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := NewInventory()
			err := inv.AddGold(tt.amount)

			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedGold, inv.GetCurrentGold())
			} else {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidGoldAmount, err)
				assert.Equal(t, 0, inv.GetCurrentGold())
			}
		})
	}
}

func TestWithdrawGold(t *testing.T) {
	tests := []struct {
		name           string
		initialGold    int
		withdrawAmount int
		shouldSucceed  bool
		expectedGold   int
	}{
		{"withdraw less than available", 100, 50, true, 50},
		{"withdraw exact amount", 100, 100, true, 0},
		{"withdraw more than available", 100, 150, false, 100},
		{"withdraw negative amount", 100, -50, false, 100},
		{"withdraw from empty inventory", 0, 50, false, 0},
		{"withdraw zero amount", 100, 0, true, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := NewInventory()
			if tt.initialGold > 0 {
				err := inv.AddGold(tt.initialGold)
				assert.NoError(t, err)
			}

			err := inv.WithdrawGold(tt.withdrawAmount)

			if tt.shouldSucceed {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedGold, inv.GetCurrentGold())
			} else {
				if tt.withdrawAmount < 0 {
					assert.Error(t, err)
					assert.Equal(t, ErrInvalidGoldAmount, err)
				} else {
					assert.Error(t, err)
					assert.Equal(t, ErrNotEnoughGold, err)
				}
				assert.Equal(t, tt.initialGold, inv.GetCurrentGold())
			}
		})
	}
}

func TestShowItems(t *testing.T) {
	tests := []struct {
		name          string
		setupItems    int
		expectedCount int
	}{
		{"show empty inventory", 0, 0},
		{"show inventory with one item", 1, 1},
		{"show inventory with multiple items", 5, 5},
		{"show inventory with max items", MAX_ITEMS, MAX_ITEMS},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := NewInventory()

			// Add items
			for i := 0; i < tt.setupItems; i++ {
				itemID := item.New(uuid.New())
				playerItem, err := playeritem.NewPlayerItem(itemID, "Test Item", 1)
				assert.NoError(t, err)
				err = inv.AddItem(*playerItem)
				assert.NoError(t, err)
			}

			items := inv.ShowItems()
			assert.Equal(t, tt.expectedCount, len(items))

			// Verify that all returned items exist in the inventory
			for _, itemID := range items {
				assert.Contains(t, inv.items, itemID)
			}
		})
	}
}

func setInventoryFull() {

}
