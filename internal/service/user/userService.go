package user_service

import (
	"context"
	"log/slog"

	user_model "github.com/Heatdog/VK-Go-PHP/internal/models/user"
	user_repository "github.com/Heatdog/VK-Go-PHP/internal/repository/user"
	cryptohash "github.com/Heatdog/VK-Go-PHP/pkg/crypto"
	"github.com/google/uuid"
)

type UserService interface {
	SignUp(context context.Context, user *user_model.UserLogin) (uuid.UUID, error)
}

type userService struct {
	logger *slog.Logger
	repo   user_repository.UserRepository
}

func NewUserService(logger *slog.Logger, repo user_repository.UserRepository) UserService {
	return &userService{
		logger: logger,
		repo:   repo,
	}
}

func (service *userService) SignUp(context context.Context, user *user_model.UserLogin) (uuid.UUID, error) {
	service.logger.Info("user service sign up")

	service.logger.Debug("hash password")
	hash, err := cryptohash.Hash(user.Password)
	if err != nil {
		service.logger.Error("hash password error")
		return uuid.Nil, nil
	}

	user.Password = string(hash)
	return service.repo.SignUp(context, user)
}
