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

func (repo *userRepositoryPostgre) InsertUser(ctx context.Context, user *user_model.UserLogin) (uuid.UUID, error) {
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

func (repo *userRepositoryPostgre) FindUser(ctx context.Context, login string) (*user_model.User, error) {
	repo.logger.Info("find user in repo")
	q := `
			SELECT id, login, password
			FROM users
			WHERE login = $1
	`
	repo.logger.Debug(q)
	row := repo.dbClient.QueryRow(ctx, q, login)

	var res user_model.User

	if err := row.Scan(&res.ID, &res.Login, &res.Password); err != nil {
		repo.logger.Error("SQL error", slog.Any("error", err))
		return nil, err
	}

	return &res, nil
}
