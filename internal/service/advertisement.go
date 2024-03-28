package service

import (
	"context"
	"fmt"
	"github.com/hanoys/marketplace-api/internal/domain"
	"image"
	_ "image/jpeg"
	"net/http"
	"time"
)

type AdvertisementServiceConfig struct {
	AdPerPage             int
	CheckImageIdleTimeout int
	MinImageWidth         int
	MaxImageWidth         int
	MinImageHeight        int
	MaxImageHeight        int
}

type AdvertisementService struct {
	repositories *Repositories
	cfg          *AdvertisementServiceConfig
}

func NewAdvertisementService(repositories *Repositories, cfg *AdvertisementServiceConfig) *AdvertisementService {
	return &AdvertisementService{repositories: repositories, cfg: cfg}
}

func checkImage(imageURL string) error {
	tr := &http.Transport{
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Transport: tr}

	res, err := client.Get(imageURL)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	imgcfg, _, err := image.DecodeConfig(res.Body)
	if err != nil {
		return fmt.Errorf("unknown format (possible formats: jpeg)")
	}

	if imgcfg.Width >= 128 && imgcfg.Height >= 128 && imgcfg.Width <= 4096 && imgcfg.Height <= 4096 {
		return nil
	}

	return fmt.Errorf("image too small (width and heigth must be at least 512p and not larger than 4096p)")
}

func (a *AdvertisementService) Create(ctx context.Context, params AdvertisementCreateParams) (domain.Advertisement, error) {
	if err := checkImage(params.ImageURL); err != nil {
		return domain.Advertisement{}, err
	}

	return a.repositories.AdvertisementsRepository.Create(ctx, params)
}

func (a *AdvertisementService) GetAdvertisements(ctx context.Context, params AdvertisementSortParams) ([]domain.AdvertisementEntry, error) {
	params.AdPerPage = a.cfg.AdPerPage
	return a.repositories.AdvertisementsRepository.GetAdvertisements(ctx, params)
}
