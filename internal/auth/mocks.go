package auth

import "github.com/swarajroy/gophersocial/internal/store"

type TestAuthenticator struct {
}

func NewTestAuthenticator() TokenGenerator {
	return &TestAuthenticator{}
}

func (ta *TestAuthenticator) GenerateToken(user *store.User) (string, error) {
	return "", nil
}

func (ta *TestAuthenticator) ValidateToken(token string) (int64, error) {
	return 0, nil
}
