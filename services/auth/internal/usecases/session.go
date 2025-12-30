package usecases

import (
	"context"
	"time"

	"github.com/ritchieridanko/erteku/services/auth/internal/infra/database"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/erteku/services/auth/internal/models"
	"github.com/ritchieridanko/erteku/services/auth/internal/repositories"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/jwt"
	"go.opentelemetry.io/otel"
)

type SessionUsecase interface {
	CreateSession(ctx context.Context, a *models.Auth, data *models.RequestMeta) (at *models.AuthToken, err *ce.Error)
	RevokeSession(ctx context.Context, refreshToken string) (err *ce.Error)
}

type sessionUsecase struct {
	appName    string
	duration   time.Duration
	sr         repositories.SessionRepository
	transactor *database.Transactor
	jwt        *jwt.JWT
}

func NewSessionUsecase(
	appName string,
	dn time.Duration,
	sr repositories.SessionRepository,
	tx *database.Transactor,
	j *jwt.JWT,
) SessionUsecase {
	return &sessionUsecase{
		appName:    appName,
		duration:   dn,
		sr:         sr,
		transactor: tx,
		jwt:        j,
	}
}

func (u *sessionUsecase) CreateSession(ctx context.Context, a *models.Auth, data *models.RequestMeta) (*models.AuthToken, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "session.usecase.CreateSession")
	defer span.End()

	// Generate refresh and access tokens
	now := time.Now().UTC()
	refreshToken := utils.GenerateUUID().String()
	accessToken, eg := u.jwt.Generate(a.ID, a.IsEmailVerified(), &now)
	if eg != nil {
		return nil, ce.NewError(
			ce.CodeJWTGenerationFailed,
			ce.MsgInternalServer,
			eg,
			logger.NewField("auth_id", a.ID),
		)
	}

	cs := models.CreateSession{
		RefreshToken: refreshToken,
		UserAgent:    data.UserAgent,
		IPAddress:    data.IPAddress,
		ExpiresAt:    now.Add(u.duration),
	}

	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		sessionID, err := u.sr.RevokeActiveSession(ctx, a.ID, data)
		if err != nil {
			return err
		}

		invalidID := int64(0)
		if sessionID != invalidID {
			cs.ParentID = &sessionID
		}

		return u.sr.CreateSession(ctx, a.ID, &cs)
	})
	return &models.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(u.duration.Seconds()),
	}, err
}

func (u *sessionUsecase) RevokeSession(ctx context.Context, refreshToken string) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "session.usecase.RevokeSession")
	defer span.End()

	// Invalid session does not fail RevokeSession process
	err := u.sr.RevokeSessionByToken(ctx, refreshToken)
	if err != nil && err.Code() != ce.CodeSessionNotFound {
		return err
	}

	return nil
}
