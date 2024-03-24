package advert_repository

import (
	"context"

	advert_model "github.com/Heatdog/VK-Go-PHP/internal/models/advert"
	"github.com/google/uuid"
)

type AdvertRepository interface {
	AddAdvert(ctx context.Context, advert *advert_model.AdvertInput,
		userID uuid.UUID) (id uuid.UUID, err error)
}
