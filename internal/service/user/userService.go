package user_service

import (
	"context"
	"log/slog"

	user_model "github.com/Heatdog/VK-Go-PHP/internal/models/user"
	"github.com/google/uuid"
)

type UserService interface {
	SignUp(context context.Context, user *user_model.UserLogin) (uuid.UUID, error)
}

type userService struct {
	logger *slog.Logger
}

func NewUserService(logger *slog.Logger) UserService {
	return &userService{
		logger: logger,
	}
}

func (service *userService) SignUp(context context.Context, user *user_model.UserLogin) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
