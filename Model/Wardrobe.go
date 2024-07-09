package Model

import "errors"

type Wardrobe struct {
	Title       string `json:"title"`
	Quantity    string `json:"quantity"`
	Price       string `json:"price"`
	Description string `json:"description"`
	Height      string `json:"height"`
	Width       string `json:"width"`
	Depth       string `json:"depth"`
	Filename    string `json:"filename"`
	Link        string `json:"link"`
}

func (w *Wardrobe) Create() error {
	if w.Title == "" {
		return errors.New("title is empty")
	}
	if w.Quantity == "" {
		return errors.New("quantity is empty")
	}
	if w.Description == "" {
		return errors.New("description is empty")
	}
	if w.Height == "" {
		return errors.New("height is empty")
	}
	if w.Width == "" {
		return errors.New("width is empty")
	}
	if w.Depth == "" {
		return errors.New("depth is empty")
	}
	return nil
}
