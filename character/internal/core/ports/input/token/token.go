package token

type AuthService interface {
	IsUserValid(userid string) bool
}
