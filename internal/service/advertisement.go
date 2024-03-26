package service

import (
	"context"
	"github.com/hanoys/marketplace-api/internal/domain"
)

type AdvertisementService struct {
	repositories *domain.Repositories
}

func NewAdvertisementService(repositories *domain.Repositories) *AdvertisementService {
	return &AdvertisementService{repositories: repositories}
}

func (a *AdvertisementService) Create(ctx context.Context, userID int, title string, body string, imageURL string, price float64) (domain.Advertisement, error) {
	return a.repositories.Advertisements.Create(ctx, userID, title, body, imageURL, price)
}
func (a *AdvertisementService) FindAll(ctx context.Context) ([]domain.Advertisement, error) {
	return nil, nil
}
