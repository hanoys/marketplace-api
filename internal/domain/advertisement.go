package domain

type Advertisement struct {
	ID       int     `json:"id"`
	UserID   int     `json:"user_id"`
	Title    string  `json:"title"`
	Body     string  `json:"body"`
	ImageURL string  `json:"image_url"`
	Price    float64 `json:"price"`
}
