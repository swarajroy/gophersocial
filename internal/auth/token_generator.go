package auth

import (
	"github.com/swarajroy/gophersocial/internal/store"
)

type AuthenticatedUser struct {
	UserID int64
}
type TokenGenerator interface {
	GenerateToken(user *store.User) (string, error)
	ValidateToken(token string) (int64, error)
}
