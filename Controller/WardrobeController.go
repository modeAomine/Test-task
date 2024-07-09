package Controller

import (
	"encoding/json"
	"net/http"
	"tests/Model"
)

type WardrobeRequest struct {
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

func AddWardrobeHandler(w http.ResponseWriter, r *http.Request) {
	var req WardrobeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wardrobe := Model.Wardrobe{
		Title:       req.Title,
		Quantity:    req.Quantity,
		Price:       req.Price,
		Description: req.Description,
		Height:      req.Height,
		Width:       req.Width,
		Depth:       req.Depth,
		Filename:    req.Filename,
		Link:        req.Link,
	}

	err = wardrobe.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": wardrobe.Title + " успешно добавлен!",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
