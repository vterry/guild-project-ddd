package login

import (
	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/base"
)

type LoginID struct {
	base.BaseID[uuid.UUID]
}

func NewLoginID(value uuid.UUID) LoginID {
	return LoginID{
		BaseID: base.New(value),
	}
}

func (l LoginID) Equals(otherID LoginID) bool {
	return l.BaseID.Equals(otherID.BaseID)
}
