package postgres

import (
	"github.com/hanoys/marketplace-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRepositories(db *pgxpool.Pool) *domain.Repositories {
	return &domain.Repositories{Users: NewUsersRepository(db),
		Advertisements: NewAdvertisementRepository(db)}
}
