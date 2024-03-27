package service

import (
	"context"
	"github.com/hanoys/marketplace-api/auth"
	"github.com/hanoys/marketplace-api/internal/domain"
)

type UsersService interface {
	SignUp(ctx context.Context, login string, password string) (domain.User, error)
}

type AdvertisementsService interface {
	Create(ctx context.Context, userID int, title string, body string, imageURL string, price float64) (domain.Advertisement, error)
	GetAdvertisements(ctx context.Context, userID, pageNumber int, sort domain.SortType, dir domain.DirectionType) ([]domain.AdvertisementEntry, error)
	FindAll(ctx context.Context) ([]domain.Advertisement, error)
}

// TODO: make token pair to be interface
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

func NewServices(repositories *Repositories, tokenProvider *auth.Provider) *Services {
	return &Services{UsersService: NewUserService(repositories),
		AdvertisementsService: NewAdvertisementService(repositories),
		Authorization:         NewAuthorizationService(repositories, tokenProvider)}
}
