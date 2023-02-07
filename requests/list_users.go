package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const defaultLimit = 10

type ListUsersRequest struct {
	Page  *uint `query:"page"`
	Limit *uint `query:"limit"`

	ID        *string `query:"id"`
	FirstName *string `query:"first_name"`
	LastName  *string `query:"last_name"`
	NickName  *string `query:"nickname" `
	Email     *string `query:"email" `
	Country   *string `query:"country"`
}

// validates list users request and return validation error
func (l ListUsersRequest) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Page, validation.When(l.Page != nil, validation.By(MinPtr(1)))),
		validation.Field(&l.Limit, validation.When(l.Limit != nil, validation.By(MinPtr(1)))),
		validation.Field(&l.ID, validation.When(l.ID != nil, validation.By(ValidateUUIDPtr))),
		validation.Field(&l.FirstName, validation.When(l.FirstName != nil, validation.By(LengthPtr(1, 255)))),
		validation.Field(&l.LastName, validation.When(l.LastName != nil, validation.By(LengthPtr(1, 255)))),
		validation.Field(&l.NickName, validation.When(l.NickName != nil, validation.By(LengthPtr(1, 255)))),
		validation.Field(&l.Email, validation.When(l.Email != nil, validation.By(LengthPtr(1, 255)))),
		validation.Field(&l.Country, validation.When(l.Country != nil, validation.By(LengthPtr(1, 2))), is.UTFLetter),
	)
}

// get limit for list request OR default
func (l ListUsersRequest) GetLimit() uint {
	if l.Limit == nil {
		return defaultLimit
	}

	return *l.Limit
}

// get offset for list request or default
func (l ListUsersRequest) GetOffset() uint {
	limit := l.GetLimit()

	var page uint = 1
	if l.Page != nil && *l.Page > 1 {
		page = *l.Page
	}

	return (page - 1) * limit
}
