package token_service

import (
	"context"
	"log/slog"

	"github.com/Heatdog/VK-Go-PHP/pkg/jwt"
)

type TokenService interface {
	GenerateToken(ctx context.Context, tokenFields jwt.TokenFileds) (accessToken string, err error)
}

type tokenService struct {
	logger    *slog.Logger
	secretKey string
}

func NewTokenService(logger *slog.Logger, secretKey string) TokenService {
	return &tokenService{
		logger:    logger,
		secretKey: secretKey,
	}
}

func (service *tokenService) GenerateToken(ctx context.Context, tokenFields jwt.TokenFileds) (string, error) {
	service.logger.Info("generate token", slog.Any("user", tokenFields.ID))

	service.logger.Debug("generate access token", slog.Any("user", tokenFields.ID))
	accessToken, err := jwt.GenerateToken(tokenFields, service.secretKey)
	if err != nil {
		service.logger.Error("generate access token failed", slog.Any("error", err))
		return "", err
	}

	return accessToken, nil
}
