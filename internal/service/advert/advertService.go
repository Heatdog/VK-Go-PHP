package advert_service

import (
	"context"
	"log/slog"
	"time"

	advert_model "github.com/Heatdog/VK-Go-PHP/internal/models/advert"
	advert_repository "github.com/Heatdog/VK-Go-PHP/internal/repository/advert"
	"github.com/google/uuid"
)

type AdvertService interface {
	AddAdvert(context context.Context, advert *advert_model.AdvertInput,
		user_id uuid.UUID) (uuid.UUID, time.Time, error)
}

type advertService struct {
	logger *slog.Logger
	repo   advert_repository.AdvertRepository
}

func NewUserService(logger *slog.Logger, repo advert_repository.AdvertRepository) AdvertService {
	return &advertService{
		logger: logger,
		repo:   repo,
	}
}

func (service *advertService) AddAdvert(context context.Context, advert *advert_model.AdvertInput,
	user_id uuid.UUID) (uuid.UUID, time.Time, error) {

	service.logger.Info("advert add service")
	return service.repo.AddAdvert(context, advert, user_id)
}
