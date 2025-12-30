package handlers

import (
	"context"

	"github.com/ritchieridanko/erteku/services/auth/internal/infra/logger"
	"github.com/ritchieridanko/erteku/services/auth/internal/models"
	"github.com/ritchieridanko/erteku/services/auth/internal/usecases"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
	"github.com/ritchieridanko/erteku/shared/contract/apis/v1"
)

type AuthHandler struct {
	apis.UnimplementedAuthServiceServer
	au     usecases.AuthUsecase
	su     usecases.SessionUsecase
	logger *logger.Logger
}

func NewAuthHandler(au usecases.AuthUsecase, su usecases.SessionUsecase, l *logger.Logger) *AuthHandler {
	return &AuthHandler{au: au, su: su, logger: l}
}

func (h *AuthHandler) SignUp(ctx context.Context, req *apis.SignUpRequest) (*apis.AuthResponse, error) {
	data1 := models.SignUpRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	a, err := h.au.SignUp(ctx, &data1)
	if err != nil {
		return nil, err
	}

	resp := apis.AuthResponse{Auth: h.toAuth(a)}

	ua, ip := utils.CtxRequestMeta(ctx)
	if ua == "" || ip == "" {
		h.logger.Warn(
			ctx,
			"incomplete request meta",
			logger.NewField("user_agent", ua),
			logger.NewField("ip_address", ip),
			logger.NewField("error_code", ce.CodeInvalidRequestMeta),
		)
		return &resp, nil
	}

	data2 := models.RequestMeta{
		UserAgent: ua,
		IPAddress: ip,
	}

	at, err := h.su.CreateSession(ctx, a, &data2)
	if err != nil {
		h.logger.Warn(
			ctx,
			"failed to create session",
			logger.NewField("auth_id", a.ID),
			logger.NewField("error_code", err.Code()),
			logger.NewField("error", err),
		)
		return &resp, nil
	}

	resp.AuthToken = &apis.AuthToken{
		AccessToken:  at.AccessToken,
		RefreshToken: at.RefreshToken,
		ExpiresIn:    at.ExpiresIn,
	}
	return &resp, nil
}

func (h *AuthHandler) toAuth(a *models.Auth) *apis.Auth {
	return &apis.Auth{
		Id:            a.ID,
		Email:         a.Email,
		EmailVerified: a.IsEmailVerified(),
	}
}
