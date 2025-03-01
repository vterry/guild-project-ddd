package item

import (
	"math/rand"

	"github.com/vterry/guild-project-ddd/domain/item/valueobjects"
)

var (
	ITEMS = []string{"Sword", "Shield", "Helmet", "Boots", "Gloves", "Potion", "Ring", "Amulet", "Scroll", "Book"}
)

type Item struct {
	valueobjects.ItemID
	Name string
}

func (i *Item) GenerateRandomItem() *Item {
	return initializeItem(ITEMS[rand.Intn(10)])
}

func initializeItem(name string) *Item {
	return &Item{
		Name: name,
	}
}
