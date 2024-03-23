package user_postgre

import (
	"context"
	"log/slog"

	user_model "github.com/Heatdog/VK-Go-PHP/internal/models/user"
	user_repository "github.com/Heatdog/VK-Go-PHP/internal/repository/user"
	"github.com/Heatdog/VK-Go-PHP/pkg/client"
	"github.com/google/uuid"
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

func (repo *userRepositoryPostgre) SignUp(ctx context.Context, user *user_model.UserLogin) (uuid.UUID, error) {
	repo.logger.Info("insert user in repo")
	q := `
			INSERT INTO users (login, password)
			VALUES ($1, $2)
			RETURNING id
	`
	repo.logger.Debug(q)

	row := repo.dbClient.QueryRow(ctx, q, user.Login, user.Password)

	var id uuid.UUID

	if err := row.Scan(&id); err != nil {
		repo.logger.Error("SQL error", slog.Any("error", err))
		return uuid.Nil, err
	}

	repo.logger.Info("successful user insert", slog.String("id", id.String()))
	return id, nil
}
