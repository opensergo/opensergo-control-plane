package util

import (
	"strings"

	"github.com/google/uuid"
)

func Nonce() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
