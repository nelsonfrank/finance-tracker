package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	JwtClaimGenerator(sub uint, exp time.Duration, iss, aud string) jwt.Claims
}
