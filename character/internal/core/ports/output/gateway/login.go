package gateway

type Login interface {
	IsValid(userId string) bool
}
