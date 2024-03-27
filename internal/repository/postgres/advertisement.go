package postgres

import (
	"context"
	"fmt"
	"github.com/hanoys/marketplace-api/internal/domain"
	"github.com/hanoys/marketplace-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
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

func formWhereClause(minPrice float64, maxPrice float64) string {
	log.Println("minPrice:", minPrice, "maxPrice:", maxPrice)
	if minPrice == 0 && maxPrice == 0 {
		return ""
	} else if minPrice == 0 && maxPrice != 0 {
		return fmt.Sprintf("WHERE a.price <= %f", maxPrice)
	} else if minPrice != 0 && maxPrice == 0 {
		return fmt.Sprintf("WHERE a.price >= %f", minPrice)
	}

	return fmt.Sprintf("WHERE a.price >= %f AND a.price <= %f", minPrice, maxPrice)
}

func makeQuery(params service.AdvertisementSortParams) string {
	var queryString string
	var dirString string
	whereClause := formWhereClause(params.MinPrice, params.MaxPrice)

	if params.Dir == domain.DefaultDir || params.Dir == domain.AscDir {
		dirString = "ASC"
	} else {
		dirString = "DESC"
	}

	if params.Sort == domain.DateSort {
		queryString = fmt.Sprintf("SELECT a.title, a.body, a.image_url, a.price, a.created_at, u.id, u.login "+
			"FROM advertisements AS a JOIN users AS u ON a.user_id = u.id %s "+
			"ORDER BY created_at %s LIMIT 2 OFFSET 2*%d",
			whereClause, dirString, params.PageNumber-1)
	} else if params.Sort == domain.PriceSort {
		queryString = fmt.Sprintf("SELECT a.title, a.body, a.image_url, a.price, a.created_at, u.id, u.login "+
			"FROM advertisements AS a JOIN users AS u ON a.user_id = u.id %s "+
			"ORDER BY price %s LIMIT 2 OFFSET 2*%d",
			whereClause, dirString, params.PageNumber-1)
	} else {
		queryString = fmt.Sprintf("SELECT a.title, a.body, a.image_url, a.price, a.created_at, u.id, u.login "+
			"FROM advertisements AS a JOIN users AS u ON a.user_id = u.id %s "+
			"LIMIT 2 OFFSET 2*%d", whereClause, params.PageNumber-1)
	}

	return queryString
}

func (r *AdvertisementRepository) GetAdvertisements(ctx context.Context, params service.AdvertisementSortParams) ([]domain.AdvertisementEntry, error) {
	rows, err := r.db.Query(ctx, makeQuery(params))
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
