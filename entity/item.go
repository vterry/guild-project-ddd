package entity

import (
	"github.com/vterry/guild-project-ddd/valueobjects"
)

type Item struct {
	valueobjects.ItemID
	Name        string
	Description string
}
