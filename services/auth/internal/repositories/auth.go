package repositories

import (
	"context"

	"github.com/ritchieridanko/erteku/services/auth/internal/models"
	"github.com/ritchieridanko/erteku/services/auth/internal/repositories/caches"
	"github.com/ritchieridanko/erteku/services/auth/internal/repositories/databases"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
)

type AuthRepository interface {
	CreateAuth(ctx context.Context, data *models.CreateAuth) (auth *models.Auth, err *ce.Error)
	GetAuthByEmail(ctx context.Context, email string) (auth *models.Auth, err *ce.Error)
	IsEmailAvailable(ctx context.Context, email string) (available bool, err *ce.Error)
}

type authRepository struct {
	database databases.AuthDatabase
	cache    caches.AuthCache
}

func NewAuthRepository(adb databases.AuthDatabase, ac caches.AuthCache) AuthRepository {
	return &authRepository{database: adb, cache: ac}
}

func (r *authRepository) CreateAuth(ctx context.Context, data *models.CreateAuth) (*models.Auth, *ce.Error) {
	return r.database.Insert(ctx, data)
}

func (r *authRepository) GetAuthByEmail(ctx context.Context, email string) (*models.Auth, *ce.Error) {
	return r.database.GetByEmail(ctx, email)
}

func (r *authRepository) IsEmailAvailable(ctx context.Context, email string) (bool, *ce.Error) {
	registered, err := r.database.IsEmailRegistered(ctx, email)
	if err != nil {
		return false, err
	}
	if registered {
		return false, nil
	}

	reserved, err := r.cache.IsEmailReserved(ctx, email)
	if err != nil {
		return false, err
	}

	return !reserved, nil
}
