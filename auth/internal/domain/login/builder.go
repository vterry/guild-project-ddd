package login

import (
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
)

type LoginBuilder struct {
	Login
}

func NewLoginBuilder() *LoginBuilder {
	return &LoginBuilder{}
}

func (b *LoginBuilder) WithLoginId(loginId LoginID) *LoginBuilder {
	b.LoginID = loginId
	return b
}

func (b *LoginBuilder) WithUserId(userId valueobjects.UserID) *LoginBuilder {
	b.userId = userId
	return b
}

func (b *LoginBuilder) WithPassword(password string) *LoginBuilder {
	b.password = password
	return b
}

func (b *LoginBuilder) Build() Login {
	return Login{
		LoginID:  b.LoginID,
		userId:   b.userId,
		password: b.password,
	}
}
