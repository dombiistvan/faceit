package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	NickName  string `json:"nickname" `
	Email     string `json:"email" `
	Password  string `json:"password"`
	Country   string `json:"country"`
}

// validates create user request and return validation error
func (c CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Required, validation.Length(3, 255)),
		validation.Field(&c.LastName, validation.Required, validation.Length(3, 255)),
		validation.Field(&c.NickName, validation.Length(0, 255)),
		validation.Field(&c.Email, validation.Required, is.EmailFormat, is.Email),
		validation.Field(&c.Password, validation.Required, validation.Required, validation.By(validatePassword), validation.Length(12, 0)),
		validation.Field(&c.Country, validation.Required, validation.Length(2, 2), is.UTFLetter),
	)
}
