package service

import (
	"context"
	"fmt"
	"github.com/hanoys/marketplace-api/internal/domain"
	"github.com/hanoys/marketplace-api/internal/domain/dto"
	"image"
	_ "image/jpeg"
	"log"
	"net/http"
)

// TODO: move to config

type AdvertisementService struct {
	repositories *domain.Repositories
}

func NewAdvertisementService(repositories *domain.Repositories) *AdvertisementService {
	return &AdvertisementService{repositories: repositories}
}

func checkImage(imageURL string) error {
	res, err := http.Get(imageURL)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	imgcfg, _, err := image.DecodeConfig(res.Body)
	if err != nil {
		return fmt.Errorf("unknown format (possible formats: jpeg)")
	}

	if imgcfg.Width >= 128 && imgcfg.Height >= 128 && imgcfg.Width <= 4096 && imgcfg.Height <= 4096 {
		log.Println("IMAGE: WIDTH:", imgcfg.Width, "HEIGHT:", imgcfg.Height)
		return nil
	}

	return fmt.Errorf("image too small (width and heigth must be at least 512p and not larger than 4096p)")
}

// Create TODO: change errors
func (a *AdvertisementService) Create(ctx context.Context, userID int, title string, body string, imageURL string, price float64) (domain.Advertisement, error) {
	if err := checkImage(imageURL); err != nil {
		return domain.Advertisement{}, err
	}

	return a.repositories.Advertisements.Create(ctx, userID, title, body, imageURL, price)
}

func (a *AdvertisementService) GetAdvertisements(ctx context.Context, userID int, pageNumber int, sort domain.SortType, dir domain.DirectionType) ([]dto.AdvertisementEntryDTO, error) {
	return a.repositories.Advertisements.GetAdvertisements(ctx, userID, pageNumber, sort, dir)
}

func (a *AdvertisementService) FindAll(ctx context.Context) ([]domain.Advertisement, error) {
	return nil, nil
}
