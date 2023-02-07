package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	NickName  *string `json:"nickname"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
	Country   *string `json:"country"`
}

// validate update user request struct
func (u UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.When(u.FirstName != nil, validation.Required, validation.By(LengthPtr(3, 255)))),
		validation.Field(&u.LastName, validation.When(u.LastName != nil, validation.Required, validation.By(LengthPtr(3, 255)))),
		validation.Field(&u.NickName, validation.When(u.NickName != nil, validation.By(LengthPtr(0, 255)))),
		validation.Field(&u.Email, validation.When(u.Email != nil, validation.Required, is.EmailFormat, is.Email)),
		validation.Field(&u.Country, validation.When(u.Country != nil, validation.Required, validation.By(LengthPtr(2, 2)), is.UTFLetter)), /**/
		validation.Field(&u.Password, validation.When(u.Password != nil, validation.Required, validation.By(validatePasswordPtr), validation.By(LengthPtr(12, 0)))),
	)
}
