package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{secret, iss, aud}
}



func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.aud),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}


func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}


func (a *JWTAuthenticator) JwtClaimGenerator(sub uint, exp time.Duration, iss, aud string) jwt.Claims {
	claims := jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": iss,
		"aud": aud,
	}

	return claims
}

func (a *JWTAuthenticator) JwtTokenGenerator(sub uint, exp time.Duration, iss, aud string) (string, error) {
	claims := a.JwtClaimGenerator(sub, exp, iss, aud)

	return a.GenerateToken(claims)
}

func (a *JWTAuthenticator) GetSubFromJWTToken(token *jwt.Token) int64 {
	claims, _ := token.Claims.(jwt.MapClaims)

	userID, _:= strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)

	return userID
}

func (a *JWTAuthenticator)GetExpiresTime(expiresIn time.Duration) time.Time  {
	now := time.Now()
	timeToExpires := now.Add(time.Duration(expiresIn) * time.Millisecond)

	return timeToExpires
}


