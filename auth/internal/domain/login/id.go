package login

import (
	"github.com/google/uuid"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
)

type LoginID struct {
	valueobjects.BaseID[uuid.UUID]
}

func NewLoginID(value uuid.UUID) LoginID {
	return LoginID{
		BaseID: valueobjects.NewBaseID(value),
	}
}

func (s LoginID) Equals(otherID LoginID) bool {
	return s.BaseID.Equals(otherID.BaseID)
}
