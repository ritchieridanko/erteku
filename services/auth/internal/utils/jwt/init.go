package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	issuer   string
	secret   string
	duration time.Duration
}

func Init(issuer, secret string, dn time.Duration) *JWT {
	return &JWT{issuer: issuer, secret: secret, duration: dn}
}

func (j *JWT) Generate(authID int64, emailVerified bool, now *time.Time) (string, error) {
	if now == nil {
		t := time.Now().UTC()
		now = &t
	}

	c := claim{
		AuthID:        authID,
		EmailVerified: emailVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   fmt.Sprintf("auth_%d", authID),
			IssuedAt:  &jwt.NumericDate{Time: *now},
			ExpiresAt: &jwt.NumericDate{Time: now.Add(j.duration)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(j.secret))
}
