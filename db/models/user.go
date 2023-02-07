package models

import (
	"crypto/sha256"
	"encoding/hex"
	"faceit/requests"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const passwordSalt string = `daoFy0jZbg52ThYToZhStCeo6HG8sliY`

type User struct {
	ID        uint   `gorm:"primaryKey" json:"-"`
	UUID      string `json:"id" gorm:"size:36;uniqueIndex"`
	FirstName string `json:"first_name" gorm:"size:255"`
	LastName  string `json:"last_name" gorm:"size:255"`
	NickName  string `json:"nickname" gorm:"size:255"`
	Email     string `json:"email" gorm:"size:255;uniqueIndex"`
	Password  string `json:"-" gorm:"size:255"`
	Country   string `json:"country" gorm:"size:2"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUserFromReq create model to insert to the database from request struct
func CreateUserFromReq(request requests.CreateUserRequest) (User, error) {
	nuuid, err := uuid.NewRandom()
	if err != nil {
		return User{}, err
	}

	pw, err := encryptPassword(request.Password)
	if err != nil {
		return User{}, fmt.Errorf("could not encrypt password: %w", err)
	}

	return User{
		UUID:      nuuid.String(),
		Email:     request.Email,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		NickName:  request.NickName,
		Country:   request.Country,
		Password:  pw,
	}, nil
}

// CreateUpdateUserFromReq creates user struct to update based on original model and update request
func CreateUpdateUserFromReq(model User, req requests.UpdateUserRequest) (User, error) {
	var err error

	pw := model.Password
	if req.Password != nil {
		pw, err = encryptPassword(*req.Password)
		if err != nil {
			return User{}, fmt.Errorf("could not encrypt password: %w", err)
		}
	}

	return User{
		ID:        model.ID,
		UUID:      model.UUID,
		Password:  pw,
		Email:     GetStringToPtrStringValue(model.Email, req.Email).(string),
		FirstName: GetStringToPtrStringValue(model.FirstName, req.FirstName).(string),
		LastName:  GetStringToPtrStringValue(model.LastName, req.LastName).(string),
		Country:   GetStringToPtrStringValue(model.Country, req.Country).(string),
		NickName:  GetStringToPtrStringValue(model.NickName, req.NickName).(string),
		CreatedAt: model.CreatedAt,
	}, nil
}

// GetStringToPtrStringValue returns string if property pointer is string pointer and original string if request prop is empty
func GetStringToPtrStringValue(origValue any, newValue any) any {
	if newValue != nil && newValue.(*string) != nil {
		return *newValue.(*string)
	}

	return origValue
}

// encryptPassword encrypts password string to hash
func encryptPassword(password string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(password + passwordSalt))
	hashPassword := hex.EncodeToString(hash.Sum(nil))

	return hashPassword, nil
}
