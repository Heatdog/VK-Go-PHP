package advert_postgre

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	advert_model "github.com/Heatdog/VK-Go-PHP/internal/models/advert"
	advert_repository "github.com/Heatdog/VK-Go-PHP/internal/repository/advert"
	"github.com/Heatdog/VK-Go-PHP/pkg/client"
	"github.com/google/uuid"
)

type advertRepositoryPostgre struct {
	dbClient client.Client
	logger   *slog.Logger
}

func NewAdvertPostgreRepository(dbClient client.Client, logger *slog.Logger) advert_repository.AdvertRepository {
	return &advertRepositoryPostgre{
		dbClient: dbClient,
		logger:   logger,
	}
}

func (repo *advertRepositoryPostgre) AddAdvert(ctx context.Context, advert *advert_model.AdvertInput,
	userID uuid.UUID) (uuid.UUID, error) {

	repo.logger.Info("add advert in repo")
	q := `
			INSERT INTO adverts (title, body, price, image_adr, user_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
	`
	repo.logger.Debug(q)

	repo.logger.Debug("params", slog.Any("advert", advert), slog.Any("user_id", userID))
	row := repo.dbClient.QueryRow(ctx, q, advert.Title, advert.Body, advert.Price, advert.ImgAddr, userID)

	var id uuid.UUID
	var dateTime time.Time

	if err := row.Scan(&id, &dateTime); err != nil {
		repo.logger.Error("SQL error", slog.Any("error", err))
		return uuid.Nil, err
	}

	repo.logger.Info("successful advert add", slog.String("id", id.String()))
	return id, nil
}

func (repo *advertRepositoryPostgre) GetAdverts(ctx context.Context,
	params advert_model.QueryParams) ([]advert_model.Advert, error) {

	repo.logger.Info("get films from repo")
	q := fmt.Sprintf(`
		SELECT a.id, a.title, a.body, a.image_adr, a.price, u.login
		FROM adverts a
		LEFT JOIN users u ON a.users_id = u.id
		WHERE a.price >= %s AND a.price <= %s
		ORDER BY %s %s

	`, params.MinPrice, params.MaxPrice, params.Sort, params.SortDir)
}
