package postgres

import (
	"context"
	"github.com/hanoys/marketplace-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdvertisementRepository struct {
	db *pgxpool.Pool
}

// TODO: public?
func NewAdvertisementRepository(db *pgxpool.Pool) *AdvertisementRepository {
	return &AdvertisementRepository{db}
}

// TODO: public?
func (r *AdvertisementRepository) Create(ctx context.Context, userID int, title string, body string, imageURL string, price float64) (domain.Advertisement, error) {
	var newAdvertisement domain.Advertisement
	err := r.db.QueryRow(ctx,
		"INSERT INTO advertisements(user_id, title, body, image_url, price) VALUES ($1, $2, $3, $4, $5) RETURNING *",
		userID,
		title,
		body,
		imageURL,
		price).Scan(
		&newAdvertisement.ID,
		&newAdvertisement.UserID,
		&newAdvertisement.Title,
		&newAdvertisement.Body,
		&newAdvertisement.ImageURL,
		&newAdvertisement.Price)

	if err != nil {
		return domain.Advertisement{}, err
	}

	return newAdvertisement, nil
}

// TODO: public?
func (r *AdvertisementRepository) FindAll(ctx context.Context) ([]domain.Advertisement, error) {
	return nil, nil
}
