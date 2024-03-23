package user_repository

import (
	"context"

	user_model "github.com/Heatdog/VK-Go-PHP/internal/models/user"
	"github.com/google/uuid"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *user_model.UserLogin) (id uuid.UUID, err error)
	FindUser(ctx context.Context, login string) (*user_model.User, error)
}
