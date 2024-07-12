package Model

type Wardrobe struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Quantity    string `json:"quantity"`
	Price       string `json:"price"`
	OldPrice    string `json:"old_price"`
	Description string `json:"description"`
	Height      string `json:"height"`
	Width       string `json:"width"`
	Depth       string `json:"depth"`
	Filename    string `json:"filename"`
	Link        string `json:"link"`
}
