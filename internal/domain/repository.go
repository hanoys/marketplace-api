package domain

import "context"

type Users interface {
	Create(ctx context.Context, login string, password string) (User, error)
}

type Advertisements interface {
	Create(ctx context.Context, userID int, title string, body string, imageURL string, price float64) (Advertisement, error)
	FindAll(ctx context.Context) ([]Advertisement, error)
}

type Repositories struct {
	Users
	Advertisements
}
