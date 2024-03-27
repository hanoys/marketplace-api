package domain

import (
	"context"
	"github.com/hanoys/marketplace-api/internal/domain/dto"
)

type Users interface {
	Create(ctx context.Context, login string, password string) (User, error)
	FindByLogin(ctx context.Context, login string) (User, error)
}

type Advertisements interface {
	Create(ctx context.Context, userID int, title string, body string, imageURL string, price float64) (Advertisement, error)
	GetAdvertisements(ctx context.Context, userID int, pageNum int, sort SortType, dir DirectionType) ([]dto.AdvertisementEntryDTO, error)
	FindAll(ctx context.Context) ([]Advertisement, error)
}

type Repositories struct {
	Users
	Advertisements
}
