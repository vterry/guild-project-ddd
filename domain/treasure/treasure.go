package treasure

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vterry/guild-project-ddd/domain/player"
	"github.com/vterry/guild-project-ddd/domain/treasure/valueobjects"
)

var (
	ErrNotEnoughCash = errors.New("player doesnt have sufficient cash")
)

// TODO Estudar como eventos de dominio poderiam ajudar aqui.
type Treasure struct {
	valueobjects.TreasureID
	cashAmount int
	donations  []*Donation
}

func NewTreasure() *Treasure {
	treasure := Treasure{
		TreasureID: valueobjects.NewTreasureID(uuid.New()),
		cashAmount: 0,
		donations:  make([]*Donation, 0),
	}

	return &treasure
}

func (t Treasure) Donate(p *player.Player, cash int) (*Treasure, error) {
	if p.GetCurrentCash() == 0 || p.GetCurrentCash() < cash {
		return nil, ErrNotEnoughCash
	}
	p.UpdateCash(-cash)
	donation := NewDonation(p.PlayerID, t.TreasureID, cash)
	t.donations = append(t.donations, &donation)
	return &t, nil
}
