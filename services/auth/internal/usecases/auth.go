package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/ritchieridanko/erteku/services/auth/internal/constants"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/database"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra/publisher"
	"github.com/ritchieridanko/erteku/services/auth/internal/models"
	"github.com/ritchieridanko/erteku/services/auth/internal/repositories"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/bcrypt"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/validator"
	"github.com/ritchieridanko/erteku/shared/contract/events/v1"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, data *models.SignUpRequest) (auth *models.Auth, err *ce.Error)
}

type authUsecase struct {
	appName    string
	ar         repositories.AuthRepository
	tr         repositories.TokenRepository
	transactor *database.Transactor
	acp        *publisher.Publisher
	validator  *validator.Validator
	bcrypt     *bcrypt.BCrypt
	logger     *logger.Logger
}

func NewAuthUsecase(
	appName string,
	ar repositories.AuthRepository,
	tr repositories.TokenRepository,
	tx *database.Transactor,
	acp *publisher.Publisher,
	v *validator.Validator,
	b *bcrypt.BCrypt,
	l *logger.Logger,
) AuthUsecase {
	return &authUsecase{
		appName:    appName,
		ar:         ar,
		tr:         tr,
		transactor: tx,
		acp:        acp,
		validator:  v,
		bcrypt:     b,
		logger:     l,
	}
}

func (u *authUsecase) SignUp(ctx context.Context, data *models.SignUpRequest) (*models.Auth, *ce.Error) {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "auth.usecase.SignUp")
	defer span.End()

	email := utils.NormalizeString(data.Email)

	// Validations
	if ok, why := u.validator.Email(&email); !ok {
		return nil, ce.NewError(ce.CodeInvalidEmail, why, nil)
	}
	if ok, why := u.validator.Password(&data.Password); !ok {
		return nil, ce.NewError(ce.CodeInvalidPassword, why, nil)
	}

	var a *models.Auth
	err := u.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		available, err := u.ar.IsEmailAvailable(ctx, email)
		if err != nil {
			return err
		}
		if !available {
			return ce.NewError(ce.CodeEmailNotAvailable, "Email is already registered", nil)
		}

		hash, eh := u.bcrypt.Hash(data.Password)
		if eh != nil {
			return ce.NewError(ce.CodeBCryptHashingFailed, ce.MsgInternalServer, eh)
		}

		data := models.CreateAuth{
			Email:    email,
			Password: hash,
		}

		a, err = u.ar.CreateAuth(ctx, &data)
		return err
	})
	if err != nil {
		return nil, err
	}

	// Store verification token
	token := utils.GenerateUUID().String()
	if err := u.tr.CreateVerification(ctx, a.ID, token); err != nil {
		u.logger.Warn(
			ctx,
			"failed to create verification token",
			logger.NewField("auth_id", a.ID),
			logger.NewField("error_code", err.Code()),
			logger.NewField("error", err),
		)
		return a, nil
	}

	// Publish auth.created event
	// fail to publish does not fail SignUp process
	ep := u.acp.Publish(
		ctx,
		fmt.Sprintf("auth_%d", a.ID),
		&events.AuthCreated{
			EventId:           utils.GenerateUUID().String(),
			Email:             email,
			VerificationToken: token,
			CreatedAt:         timestamppb.New(time.Now().UTC()),
		},
	)
	if ep != nil {
		u.logger.Warn(
			ctx,
			"failed to publish event",
			logger.NewField("event_topic", constants.EventTopicAC),
			logger.NewField("auth_id", a.ID),
			logger.NewField("error_code", ce.CodeEventPublishingFailed),
			logger.NewField("error", ep),
		)
	}

	return a, nil
}
