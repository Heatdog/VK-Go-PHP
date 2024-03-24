package advert_service

import (
	"context"
	"log/slog"

	advert_model "github.com/Heatdog/VK-Go-PHP/internal/models/advert"
	advert_repository "github.com/Heatdog/VK-Go-PHP/internal/repository/advert"
	"github.com/google/uuid"
)

type AdvertService interface {
	AddAdvert(context context.Context, advert *advert_model.AdvertInput,
		user_id uuid.UUID) (uuid.UUID, error)

	GetAdverts(context context.Context, params advert_model.QueryParams,
		userID uuid.UUID) ([]advert_model.Advert, error)
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
	user_id uuid.UUID) (uuid.UUID, error) {

	service.logger.Info("advert add service")
	return service.repo.AddAdvert(context, advert, user_id)
}

func (service *advertService) GetAdverts(context context.Context, params advert_model.QueryParams,
	userID uuid.UUID) ([]advert_model.Advert, error) {

	service.logger.Info("get adverts info")
	adverts, err := service.repo.GetAdverts(context, params)
	if err != nil {
		service.logger.Warn(err.Error())
		return nil, err
	}

	for i := range adverts {
		if adverts[i].UserID == userID {
			adverts[i].Own = true
		} else {
			adverts[i].Own = false
		}
	}

	return adverts, nil
}
