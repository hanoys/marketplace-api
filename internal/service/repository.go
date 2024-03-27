package service

import (
	"context"
	"github.com/hanoys/marketplace-api/internal/domain"
)

type UsersRepository interface {
	Create(ctx context.Context, login string, password string) (domain.User, error)
	FindByLogin(ctx context.Context, login string) (domain.User, error)
}

type AdvertisementsRepository interface {
	Create(ctx context.Context, params AdvertisementCreateParams) (domain.Advertisement, error)
	GetAdvertisements(ctx context.Context, params AdvertisementSortParams) ([]domain.AdvertisementEntry, error)
}

type Repositories struct {
	UsersRepository
	AdvertisementsRepository
}
