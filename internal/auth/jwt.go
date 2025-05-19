package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/swarajroy/gophersocial/internal/store"
)

type JWTTokenGenerator struct {
	secret string
	host   string
	exp    time.Duration
}

func NewJWTTokenGenerator(secret, host string, exp time.Duration) *JWTTokenGenerator {
	return &JWTTokenGenerator{
		secret: secret,
		host:   host,
		exp:    exp,
	}
}

func (jtg *JWTTokenGenerator) GenerateToken(user *store.User) (string, error) {
	claims := jwt.MapClaims{
		"iss": jtg.host,
		"sub": user.ID,
		"aud": jtg.host,
		"exp": time.Now().Add(jtg.exp),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jtg.secret)) // why slice of bytes???
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (jtg *JWTTokenGenerator) ValidateToken(token string) (string, error) {
	return "", nil
}
