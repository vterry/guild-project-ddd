package login

import "github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"

type Repository interface {
	Save(login Login) error
	FindLoginByUserID(userId valueobjects.UserID) (*Login, error)
	UpdatePassword(userId valueobjects.UserID, password string) error
}
