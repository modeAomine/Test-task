package Controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tests/Model"
	"tests/Service"
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

	err = Service.CreateWardrobe(&wardrobe)
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

func UpdateWardrobeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var wardrobe Model.Wardrobe
	err = json.NewDecoder(r.Body).Decode(&wardrobe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wardrobe.ID = id

	err = Service.UpdateWardrobe(&wardrobe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": wardrobe.Title + " успешно отредактирован!",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func DeleteWardrobeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = Service.DeleteWardrobe(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Wardrobe deleted successfully!",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
