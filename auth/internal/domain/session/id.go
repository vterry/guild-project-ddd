package session

import (
	"github.com/google/uuid"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
)

type SessionID struct {
	valueobjects.BaseID[uuid.UUID]
}

func NewSessionID(value uuid.UUID) SessionID {
	return SessionID{
		BaseID: valueobjects.NewBaseID(value),
	}
}

func (s SessionID) Equals(otherID SessionID) bool {
	return s.BaseID.Equals(otherID.BaseID)
}
