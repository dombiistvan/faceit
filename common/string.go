package common

import (
	"strings"
)

func IsEmptyString(value string) bool {
	return strings.TrimSpace(value) == ``
}
