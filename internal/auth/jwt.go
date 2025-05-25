package auth

import (
	"fmt"
	"log"
	"strconv"
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
		"exp": time.Now().Add(jtg.exp).Unix(),
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

func (jtg *JWTTokenGenerator) ValidateToken(token string) (int64, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(jtg.secret), nil
	},
		jwt.WithAudience(jtg.host),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(jtg.host),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		log.Printf("error is %+v \n", err)
		return 0, err
	}
	claims, _ := jwtToken.Claims.(jwt.MapClaims)

	userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
