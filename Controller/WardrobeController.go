package Controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"tests/Model"
	"tests/Service"
)

type WardrobeRequest struct {
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

func AddWardrobeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(25 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	quantity := r.FormValue("quantity")
	price := r.FormValue("price")
	oldPrice := r.FormValue("old_price")
	description := r.FormValue("description")
	height := r.FormValue("height")
	width := r.FormValue("width")
	depth := r.FormValue("depth")
	link := r.FormValue("link")

	file, handler, err := r.FormFile("filename")
	if err != nil {
		http.Error(w, "Failed to get filename", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := handler.Filename

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		return
	}

	wardrobe := Model.Wardrobe{
		Title:       title,
		Quantity:    quantity,
		Price:       price,
		OldPrice:    oldPrice,
		Description: description,
		Height:      height,
		Width:       width,
		Depth:       depth,
		Filename:    filename,
		Link:        link,
	}

	err = Service.CreateWardrobe(&wardrobe, fileBytes)
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

func GetWardrobeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wadrobe, err := Service.GetWardrobeById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if wadrobe == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wadrobe)
}
