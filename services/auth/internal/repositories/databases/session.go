package databases

import (
	"context"
	"errors"
	"time"

	"github.com/ritchieridanko/erteku/services/auth/internal/infra/database"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/erteku/services/auth/internal/models"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
)

type SessionDatabase interface {
	Insert(ctx context.Context, authID int64, data *models.CreateSession) (err *ce.Error)
	RevokeByToken(ctx context.Context, refreshToken string) (err *ce.Error)
	RevokeActive(ctx context.Context, authID int64, data *models.RequestMeta) (sessionID int64, err *ce.Error)
}

type sessionDatabase struct {
	database *database.Database
}

func NewSessionDatabase(db *database.Database) SessionDatabase {
	return &sessionDatabase{database: db}
}

func (d *sessionDatabase) Insert(ctx context.Context, authID int64, data *models.CreateSession) *ce.Error {
	query := `
		INSERT INTO sessions (
			auth_id, parent_id, refresh_token,
			user_agent, ip_address, expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	err := d.database.Execute(
		ctx, query, authID, data.ParentID, data.RefreshToken,
		data.UserAgent, data.IPAddress, data.ExpiresAt,
	)
	if err != nil {
		return ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			err,
			logger.NewField("auth_id", authID),
		)
	}

	return nil
}

func (d *sessionDatabase) RevokeByToken(ctx context.Context, refreshToken string) *ce.Error {
	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE refresh_token = $1 AND expires_at >= $2 AND revoked_at IS NULL
	`

	if err := d.database.Execute(ctx, query, refreshToken, time.Now().UTC()); err != nil {
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(ce.CodeSessionNotFound, ce.MsgResourceNotFound, err)
		}
		return ce.NewError(ce.CodeDBQueryExec, ce.MsgInternalServer, err)
	}
	return nil
}

func (d *sessionDatabase) RevokeActive(ctx context.Context, authID int64, data *models.RequestMeta) (int64, *ce.Error) {
	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE
			auth_id = $1 AND user_agent = $2 AND ip_address = $3
			AND expires_at >= $4 AND revoked_at IS NULL
		RETURNING id
	`

	row := d.database.Query(
		ctx, query, authID, data.UserAgent,
		data.IPAddress, time.Now().UTC(),
	)

	var sessionID int64
	if err := row.Scan(&sessionID); err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return 0, nil
		}
		return 0, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			err,
			logger.NewField("auth_id", authID),
		)
	}

	return sessionID, nil
}
