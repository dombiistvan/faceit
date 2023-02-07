package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"regexp"
)

var (
	ErrAlreadyExist   = errors.New("record already exist with the value")
	ErrInternal       = errors.New("internal error")
	ErrNotFoundSimple = errors.New("entity/entities not found")
	ErrInvalidModel   = errors.New("invalid model")
)

type ValidationErrorMap map[string]interface{}

//json parse error for general usage
func ErrJsonParse(err error) error {
	return fmt.Errorf("could not parse json data: %w", err)
}

// return json marshall error for general usage
func ErrMarshall(err error) error {
	return fmt.Errorf("could not marshall: %w", err)
}

// return json unmarshall error for general usage
func ErrUnmarshall(object string, err error) error {
	return fmt.Errorf("could not unmarshall %s: %w", object, err)
}

// return bind request error for general usage
func ErrBindRequest(err error) error {
	return fmt.Errorf("could not bind request: %w", err)
}

// return not found error for general usage
func ErrNotFound(model string, err error) error {
	return fmt.Errorf("could not find %s: %w", model, err)
}

// ValidationErrorMap is a wrapper for proper display of ozzo validation error after validated struct
func (vem ValidationErrorMap) Error() string {
	marshalledErr, err := json.Marshal(vem)
	if err != nil {
		return ErrMarshall(err).Error()
	}

	return string(marshalledErr)
}

// marshall object
func (vem ValidationErrorMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(vem))
}

// password validation function of requests
func validatePassword(value interface{}) error {
	var regexpList = []string{`[a-z]+`, `[A-Z]+`, `[\!\@\#\$\%\&\*\+\-\_\=\?\:\;\,\.\|\(\)\{\}\<\>]+`, `[0-9]+`}
	for _, pattern := range regexpList {
		var re = regexp.MustCompile(pattern)
		if !re.MatchString(value.(string)) {
			return fmt.Errorf(`the password has to contain at least one lowercase and uppercase letter, a special character (\!\@\#\$\%%\&\*\+\-\_\=\?\:\;\,\.\|\(\)\{\}\<\>) and a number`)
		}
	}

	return nil
}

// password validation function to pointers
func validatePasswordPtr(value interface{}) error {
	return validatePassword(*value.(*string))
}

// length validation rule for str pointer properties
func LengthPtr(minLength, maxLength int) validation.RuleFunc {
	return func(value interface{}) error {
		if (minLength > 0 && len(*value.(*string)) < minLength) || (maxLength > 0 && len(*value.(*string)) > maxLength) {
			return fmt.Errorf("length must be between %d and %d", minLength, maxLength)
		}
		return nil
	}
}

// uuid validation rule for uuid pointer properties
func ValidateUUIDPtr(value interface{}) error {
	if value == nil {
		return nil
	}

	if value.(*string) == nil {
		return nil
	}

	_, err := uuid.Parse(*value.(*string))
	if err == nil {
		return nil
	}

	return fmt.Errorf("invalid UUID")
}

// min validation rule for int pointer properties
func MinPtr(minValue int) validation.RuleFunc {
	return func(value interface{}) error {
		var cValue int
		switch value.(type) {
		case *int:
			cValue = *value.(*int)
			break
		case *uint:
			cValue = int(*value.(*uint))
			break
		default:
			return errors.New("unprepared type for MinPtr")
		}
		if cValue < minValue {
			return fmt.Errorf("value must be minimum %d", minValue)
		}
		return nil
	}
}
