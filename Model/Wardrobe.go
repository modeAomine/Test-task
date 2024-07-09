package Model

import (
	"errors"
	"tests/DataBase"
)

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

func (w *Wardrobe) Create() error {
	requiredFields := map[string]interface{}{
		"title":       w.Title,
		"quantity":    w.Quantity,
		"price":       w.Price,
		"description": w.Description,
		"height":      w.Height,
		"width":       w.Width,
		"depth":       w.Depth,
		"filename":    w.Filename,
		"link":        w.Link,
	}

	for field, value := range requiredFields {
		if value == "" || value == 0 {
			return errors.New(field + " is required")
		}
	}

	query := `INSERT INTO Wardrobe (title, quantity, price, description, height, width, depth, filename, link) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	err := DataBase.DB.QueryRow(query, w.Title, w.Quantity, w.Price, w.Description, w.Height, w.Width, w.Depth, w.Filename, w.Link).Scan(&w.ID)
	if err != nil {
		return err
	}

	return nil
}
