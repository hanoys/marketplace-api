package postgres

import (
	"github.com/hanoys/marketplace-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRepositories(db *pgxpool.Pool) *service.Repositories {
	return &service.Repositories{UsersRepository: NewUsersRepository(db),
		AdvertisementsRepository: NewAdvertisementRepository(db)}
}
