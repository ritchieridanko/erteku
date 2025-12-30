package repositories

import (
	"context"

	"github.com/ritchieridanko/erteku/services/auth/internal/models"
	"github.com/ritchieridanko/erteku/services/auth/internal/repositories/databases"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, authID int64, data *models.CreateSession) (err *ce.Error)
	RevokeSessionByToken(ctx context.Context, refreshToken string) (err *ce.Error)
	RevokeActiveSession(ctx context.Context, authID int64, data *models.RequestMeta) (sessionID int64, err *ce.Error)
}

type sessionRepository struct {
	database databases.SessionDatabase
}

func NewSessionRepository(sdb databases.SessionDatabase) SessionRepository {
	return &sessionRepository{database: sdb}
}

func (r *sessionRepository) CreateSession(ctx context.Context, authID int64, data *models.CreateSession) *ce.Error {
	return r.database.Insert(ctx, authID, data)
}

func (r *sessionRepository) RevokeSessionByToken(ctx context.Context, refreshToken string) *ce.Error {
	return r.database.RevokeByToken(ctx, refreshToken)
}

func (r *sessionRepository) RevokeActiveSession(ctx context.Context, authID int64, data *models.RequestMeta) (int64, *ce.Error) {
	return r.database.RevokeActive(ctx, authID, data)
}
