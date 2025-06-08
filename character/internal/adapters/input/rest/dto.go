package rest

type CreateCharacterRequest struct {
	UserID   string `json:"userId" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Nickname string `json:"nickname" validate:"required"`
	Class    string `json:"class" validate:"required"`
}
