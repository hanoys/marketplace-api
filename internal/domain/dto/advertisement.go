package dto

type PostAdvertisementDTO struct {
	UserID   int     `json:"user_id" binding:"required"`
	Title    string  `json:"title" binding:"required"`
	Body     string  `json:"body" binding:"required"`
	ImageURL string  `json:"image_url" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
}
