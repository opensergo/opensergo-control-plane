package util

import (
	"github.com/google/uuid"
	"strings"
)

func Nonce() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
