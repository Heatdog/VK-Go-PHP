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
		userID uuid.UUID) ([][]advert_model.AdvertWithOwner, error)
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
	userID uuid.UUID) ([][]advert_model.AdvertWithOwner, error) {

	service.logger.Info("get adverts info")
	adverts, err := service.repo.GetAdverts(context, params)
	if err != nil {
		service.logger.Warn(err.Error())
		return nil, err
	}
	service.logger.Debug("adverts", slog.Any("adverts", adverts))

	if userID != uuid.Nil {
		for i := range adverts {
			if adverts[i].UserID == userID {
				adverts[i].Own = true
			}
		}
	}

	var res [][]advert_model.AdvertWithOwner
	i, j := 0, 0
	for _, el := range adverts {
		if i == 0 {
			list := make([]advert_model.AdvertWithOwner, 0, 10)
			res = append(res, list)
		}

		service.logger.Debug("append", slog.Int("i", i), slog.Int("j", j), slog.Any("el", el))
		res[j] = append(res[j], el)
		i++
		if i == 10 {
			i = 0
			j++
		}
	}

	return res, nil
}
