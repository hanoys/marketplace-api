package dto

import "time"

type PostAdvertisementDTO struct {
	Title    string  `json:"title" binding:"required,max=255"`
	Body     string  `json:"body" binding:"required,max=2048"`
	ImageURL string  `json:"image_url" binding:"required,max=2048"`
	Price    float64 `json:"price" binding:"required"`
}

type AdvertisementEntryDTO struct {
	PostedByYou bool      `json:"posted_by_you"`
	UserLogin   string    `json:"user_login"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}
