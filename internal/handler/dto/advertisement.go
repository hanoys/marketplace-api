package dto

type PostAdvertisementDTO struct {
	Title    string  `json:"title" binding:"required,max=255"`
	Body     string  `json:"body" binding:"required,max=2048"`
	ImageURL string  `json:"image_url" binding:"required,max=2048"`
	Price    float64 `json:"price" binding:"required"`
}
