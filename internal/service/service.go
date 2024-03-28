package service

import (
	"context"
	"github.com/hanoys/marketplace-api/auth"
	"github.com/hanoys/marketplace-api/internal/domain"
)

type UsersService interface {
	SignUp(ctx context.Context, login string, password string) (domain.User, error)
}

type AdvertisementCreateParams struct {
	UserID   int
	Title    string
	Body     string
	ImageURL string
	Price    float64
}

type AdvertisementSortParams struct {
	UserID     int
	PageNumber int
	Sort       domain.SortType
	Dir        domain.DirectionType
	MinPrice   float64
	MaxPrice   float64
	AdPerPage  int
}

type AdvertisementsService interface {
	Create(ctx context.Context, createParams AdvertisementCreateParams) (domain.Advertisement, error)
	GetAdvertisements(ctx context.Context, sortParams AdvertisementSortParams) ([]domain.AdvertisementEntry, error)
}

type Authorization interface {
	LogIn(ctx context.Context, login string, password string) (*auth.TokenPair, error)
	LogOut(ctx context.Context, tokenString string) error
	RefreshToken(ctx context.Context, refreshTokenString string) (*auth.TokenPair, error)
	VerifyToken(ctx context.Context, tokenString string) (int, error)
}

type Services struct {
	UsersService
	AdvertisementsService
	Authorization
}

func NewServices(repositories *Repositories, tokenProvider *auth.Provider, cfg *AdvertisementServiceConfig) *Services {
	return &Services{
		UsersService:          NewUserService(repositories),
		AdvertisementsService: NewAdvertisementService(repositories, cfg),
		Authorization:         NewAuthorizationService(repositories, tokenProvider)}
}
