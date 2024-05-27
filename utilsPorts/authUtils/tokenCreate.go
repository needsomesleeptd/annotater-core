package auth_utils

import (
	"errors"

	"github.com/needsomesleeptd/annotater-core/models"
)

type Payload struct {
	Login string
	ID    uint64
	Role  models.Role
}

type ITokenHandler interface {
	GenerateToken(credentials models.User, key string) (string, error)
	ValidateToken(tokenString string, key string) error
	ParseToken(tokenString string, key string) (*Payload, error)
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrParsingToken = errors.New("error parsing token")
)
