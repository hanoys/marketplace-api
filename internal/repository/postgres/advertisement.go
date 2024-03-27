package postgres

import (
	"context"
	"fmt"
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
		"INSERT INTO advertisements(user_id, title, body, image_url, price, created_at) VALUES ($1, $2, $3, $4, $5, now()) RETURNING *",
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
		&newAdvertisement.Price,
		&newAdvertisement.CreatedAt)

	if err != nil {
		return domain.Advertisement{}, err
	}

	return newAdvertisement, nil
}

func makeQuery(pageNum int, sort domain.SortType, dir domain.DirectionType) string {
	var queryString string
	var dirString string

	if dir == domain.DefaultDir || dir == domain.AscDir {
		dirString = "ASC"
	} else {
		dirString = "DESC"
	}

	if sort == domain.DateSort {
		queryString = fmt.Sprintf("SELECT a.title, a.body, a.image_url, a.price, a.created_at, u.id, u.login "+
			"FROM advertisements AS a JOIN users AS u ON a.user_id = u.id "+
			"ORDER BY created_at %s LIMIT 2 OFFSET 2*%d",
			dirString, pageNum)
	} else if sort == domain.PriceSort {
		queryString = fmt.Sprintf("SELECT a.title, a.body, a.image_url, a.price, a.created_at, u.id, u.login "+
			"FROM advertisements AS a JOIN users AS u ON a.user_id = u.id "+
			"ORDER BY price %s LIMIT 2 OFFSET 2*%d",
			dirString, pageNum)
	} else {
		queryString = fmt.Sprintf("SELECT a.title, a.body, a.image_url, a.price, a.created_at, u.id, u.login "+
			"FROM advertisements AS a JOIN users AS u ON a.user_id = u.id "+
			"LIMIT 2 OFFSET 2*%d", pageNum)
	}

	return queryString
}

func (r *AdvertisementRepository) GetAdvertisements(ctx context.Context, userID int, pageNum int, sort domain.SortType, dir domain.DirectionType) ([]domain.AdvertisementEntry, error) {

	rows, err := r.db.Query(ctx, makeQuery(pageNum-1, sort, dir))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var advertisements []domain.AdvertisementEntry
	for rows.Next() {
		var ad domain.AdvertisementEntry
		var id int
		if err = rows.Scan(
			&ad.Title,
			&ad.Body,
			&ad.ImageURL,
			&ad.Price,
			&ad.CreatedAt,
			&id,
			&ad.UserLogin); err != nil {
			return nil, err
		}

		if id == userID {
			ad.PostedByYou = true
		}
		advertisements = append(advertisements, ad)
	}

	return advertisements, nil
}

// TODO: public?
func (r *AdvertisementRepository) FindAll(ctx context.Context) ([]domain.Advertisement, error) {
	return nil, nil
}
