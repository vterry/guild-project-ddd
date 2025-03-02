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
	name string
}

func PickRandomItem() *Item {
	return initializeItem(ITEMS[rand.Intn(10)])
}

func (i *Item) Name() string {
	return i.name
}

func initializeItem(name string) *Item {
	return &Item{
		name: name,
	}
}
