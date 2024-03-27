package domain

import "time"

type Advertisement struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	ImageURL  string    `json:"image_url"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type SortType int

const (
	DefaultSort SortType = iota
	PriceSort
	DateSort
)

type DirectionType int

const (
	DefaultDir DirectionType = iota
	AscDir
	DescDir
)
