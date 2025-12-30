package models

import "time"

type Auth struct {
	ID                int64
	Email             string
	Password          *string
	EmailVerifiedAt   *time.Time
	EmailChangedAt    *time.Time
	PasswordChangedAt *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

func (a *Auth) IsEmailVerified() bool {
	return a.EmailVerifiedAt != nil
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type CreateAuth struct {
	Email    string
	Password string
}
