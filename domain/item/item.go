package item

import "github.com/vterry/guild-project-ddd/domain/item/valueobjects"

type Item struct {
	valueobjects.ItemID
	Name        string
	Description string
}
