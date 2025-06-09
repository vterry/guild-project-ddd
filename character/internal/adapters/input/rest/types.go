package rest

type CreateCharacterRequest struct {
	UserID   string `json:"userId" validate:"required"`
	Nickname string `json:"nickname" validate:"required"`
	Class    string `json:"class" validate:"required"`
}
