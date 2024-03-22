package user_postgre

import (
	"log/slog"

	user_repository "github.com/Heatdog/VK-Go-PHP/internal/repository/user"
	"github.com/Heatdog/VK-Go-PHP/pkg/client"
)

type userRepositoryPostgre struct {
	dbClient client.Client
	logger   *slog.Logger
}

func NewUserPostgreRepository(dbClient client.Client, logger *slog.Logger) user_repository.UserRepository {
	return &userRepositoryPostgre{
		dbClient: dbClient,
		logger:   logger,
	}
}
