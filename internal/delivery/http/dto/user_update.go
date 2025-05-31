package dto

type UpdateUserRequest struct {
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty" validate:"omitempty,min=6"`
	Name     string `json:"name,omitempty" validate:"omitempty,min=3"`
	Surname  string `json:"surname,omitempty" validate:"omitempty,min=3"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,min=6"`
	Telegram string `json:"telegram,omitempty" validate:"omitempty,min=4"`
}
