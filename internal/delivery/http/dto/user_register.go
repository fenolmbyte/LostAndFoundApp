package dto

type UserRegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name"     validate:"required,min=3"`
	Surname  string `json:"surname"  validate:"required,min=3"`
	Phone    string `json:"phone"    validate:"required,min=6"`
	Telegram string `json:"telegram" validate:"required,min=4"`
}
