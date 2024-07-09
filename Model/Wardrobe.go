package Model

type Wardrobe struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Height      float64 `json:"height"`
	Width       float64 `json:"width"`
	Depth       float64 `json:"depth"`
	Filename    string  `json:"filename"`
	Link        string  `json:"link"`
}
