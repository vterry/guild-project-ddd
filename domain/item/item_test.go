package item

import (
	"testing"

	"github.com/google/uuid"
)

func TestPickRandomItem(t *testing.T) {

	t.Run("Pick a random item", func(t *testing.T) {

		item := PickRandomItem()

		if item == nil {
			t.Errorf("Expected item to be not nil")
		}

		found := false
		for _, name := range ITEMS {
			if item.Name() == name {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected item to be one of the predefined items")
		}
	})
}

func TestInitializeItem(t *testing.T) {
	name := "sword"
	item := initializeItem(name)

	if item.Name() != name {
		t.Errorf("Expected item name to be %s, but got %s", name, item.Name())
	}
}
func TestItemID(t *testing.T) {
	item := initializeItem("Sword")
	if item.ID() != uuid.Nil {
		t.Errorf("Expected item ID to be not empty")
	}
}
