package types

type AuthRegister struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUser struct {
	ID      uint
	Updates map[string]interface{}
}
