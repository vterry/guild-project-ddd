package entity

import (
	"github.com/vterry/guild-project-ddd/valueobjects"
)

type Treasure struct {
	valueobjects.TreasureID
	CashAmount int
}
