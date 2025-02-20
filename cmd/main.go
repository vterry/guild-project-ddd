package main

import (
	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/aggregate"
	"github.com/vterry/guild-project-ddd/entity"
	"github.com/vterry/guild-project-ddd/valueobjects"
)

func main() {

	// Players

	p1 := entity.NewPlayer("Player 1", valueobjects.Warrior)
	p2 := entity.NewPlayer("Player 2", valueobjects.Ranger)

	//Items

	sword := entity.Item{
		ItemID:      valueobjects.NewItemID(uuid.New()),
		Name:        "Regular Sword",
		Description: "Regular Sword - there is nothing special about it",
	}

	staff := entity.Item{
		ItemID:      valueobjects.NewItemID(uuid.New()),
		Name:        "Regular Staff",
		Description: "Same as regular Sword",
	}

	//Guild

	g1, _ := aggregate.CreateGuild("New Guild", p1)
	g1.AddPlayer(p2)

	g1.AddItem(&sword)
	g1.AddItem(&staff)

	g1.Print()
}
