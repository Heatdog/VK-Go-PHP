package user_service

import (
	"context"
	"fmt"
	"log/slog"

	user_model "github.com/Heatdog/VK-Go-PHP/internal/models/user"
	user_repository "github.com/Heatdog/VK-Go-PHP/internal/repository/user"
	token_service "github.com/Heatdog/VK-Go-PHP/internal/service/token"
	cryptohash "github.com/Heatdog/VK-Go-PHP/pkg/crypto"
	"github.com/Heatdog/VK-Go-PHP/pkg/jwt"
	"github.com/google/uuid"
)

type UserService interface {
	SignUp(context context.Context, user *user_model.UserLogin) (uuid.UUID, error)
	SignIn(context context.Context, user *user_model.UserLogin) (string, error)
}

type userService struct {
	logger       *slog.Logger
	repo         user_repository.UserRepository
	tokenService token_service.TokenService
}

func NewUserService(logger *slog.Logger, repo user_repository.UserRepository,
	tokenService token_service.TokenService) UserService {
	return &userService{
		logger:       logger,
		repo:         repo,
		tokenService: tokenService,
	}
}

func (service *userService) SignUp(context context.Context, user *user_model.UserLogin) (uuid.UUID, error) {
	service.logger.Info("user service sign up")

	service.logger.Debug("hash password")
	hash, err := cryptohash.Hash(user.Password)
	if err != nil {
		service.logger.Error(err.Error())
		return uuid.Nil, err
	}

	user.Password = string(hash)
	return service.repo.InsertUser(context, user)
}

func (service *userService) SignIn(context context.Context, user *user_model.UserLogin) (string, error) {
	service.logger.Info("user service sign up")

	service.logger.Debug("get user")
	repoUser, err := service.repo.FindUser(context, user.Login)
	if err != nil {
		service.logger.Error(err.Error())
		return "", err
	}

	service.logger.Debug("verify passwords")
	if !cryptohash.VerifyHash([]byte(repoUser.Password), user.Password) {
		service.logger.Info("passwords are not equal")
		return "", fmt.Errorf("passwords are not equal")
	}

	return service.tokenService.GenerateToken(context, jwt.TokenFileds{
		ID: repoUser.ID.String(),
	})
}
