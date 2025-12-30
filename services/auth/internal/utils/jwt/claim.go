package jwt

import "github.com/golang-jwt/jwt/v5"

type claim struct {
	AuthID        int64
	EmailVerified bool
	jwt.RegisteredClaims
}
