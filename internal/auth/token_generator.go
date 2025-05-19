package auth

import (
	"github.com/swarajroy/gophersocial/internal/store"
)

type TokenGenerator interface {
	GenerateToken(user *store.User) (string, error)
	ValidateToken(token string) (string, error)
}
