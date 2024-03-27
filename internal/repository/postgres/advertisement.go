package postgres

import (
	"context"
	"fmt"
	"github.com/hanoys/marketplace-api/internal/domain"
	"github.com/hanoys/marketplace-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdvertisementRepository struct {
	db *pgxpool.Pool
}

func NewAdvertisementRepository(db *pgxpool.Pool) *AdvertisementRepository {
	return &AdvertisementRepository{db}
}

func (r *AdvertisementRepository) Create(ctx context.Context, params service.AdvertisementCreateParams) (domain.Advertisement, error) {
	var newAdvertisement domain.Advertisement
	err := r.db.QueryRow(ctx,
		"INSERT INTO advertisements(user_id, title, body, image_url, price, created_at) VALUES ($1, $2, $3, $4, $5, now()) RETURNING *",
		params.UserID, params.Title, params.Body, params.ImageURL, params.Price).Scan(
		&newAdvertisement.ID, &newAdvertisement.UserID, &newAdvertisement.Title, &newAdvertisement.Body,
		&newAdvertisement.ImageURL, &newAdvertisement.Price, &newAdvertisement.CreatedAt)

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

func (r *AdvertisementRepository) GetAdvertisements(ctx context.Context, params service.AdvertisementSortParams) ([]domain.AdvertisementEntry, error) {

	rows, err := r.db.Query(ctx, makeQuery(params.PageNumber-1, params.Sort, params.Dir))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var advertisements []domain.AdvertisementEntry
	for rows.Next() {
		var ad domain.AdvertisementEntry
		var id int
		if err = rows.Scan(&ad.Title, &ad.Body, &ad.ImageURL, &ad.Price, &ad.CreatedAt,
			&id, &ad.UserLogin); err != nil {
			return nil, err
		}

		if id == params.UserID {
			ad.PostedByYou = true
		}
		advertisements = append(advertisements, ad)
	}

	return advertisements, nil
}
