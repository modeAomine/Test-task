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
	if !Service.CheckTokenExpiration(w, r) {
		return
	}

	err := r.ParseMultipartForm(25 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	errors := make(map[string]string)

	title := r.FormValue("title")
	quantity := r.FormValue("quantity")
	price := r.FormValue("price")
	oldPrice := r.FormValue("old_price")
	description := r.FormValue("description")
	height := r.FormValue("height")
	width := r.FormValue("width")
	depth := r.FormValue("depth")
	link := r.FormValue("link")

	if title == "" {
		errors["title"] = "Название не может быть пустым"
	}
	if quantity == "" {
		errors["quantity"] = "Количество не может быть пустым"
	}
	if price == "" {
		errors["price"] = "Цена не может быть пустой"
	}

	if oldPrice == "" {
		errors["old_price"] = "Старая цена не может быть пустой"
	}

	if oldPrice > price {
		errors["old_price"] = "Старая цена не может быть больше текущей"
	}

	if oldPrice == "" {
		errors["old_price"] = "Старая цена не может быть пустой"
	}

	if description == "" {
		errors["description"] = "Описание не может быть пустым"
	}

	if height == "" || width == "" || depth == "" {
		errors["sizes"] = "Высота, ширина и грубина не можгут быть пустыми"
	}

	if link == "" {
		errors["link"] = "Ссылка не может быть пустой"
	}

	file, handler, err := r.FormFile("filename")
	if err != nil {
		errors["filename"] = "Не удалось получить файл"
	} else {
		defer file.Close()

		filename := handler.Filename

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			errors["filename"] = "Не получилось загрузить фотографию"
		}

		if len(errors) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"errors": errors})
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

	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errors})
	}
}

func UpdateWardrobeHandler(w http.ResponseWriter, r *http.Request) {
	if !Service.CheckTokenExpiration(w, r) {
		return
	}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(25 << 20)
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

	var fileBytes []byte
	var filename string

	file, handler, err := r.FormFile("filename")
	if err == nil {
		defer file.Close()
		filename = handler.Filename
		fileBytes, err = ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Ошибка загрузки фотографии", http.StatusInternalServerError)
			return
		}
	} else if err == http.ErrMissingFile {
		existingWardrobe, err := Service.GetWardrobeById(id)
		if err != nil {
			http.Error(w, "Не удалось получить текущий шкаф", http.StatusInternalServerError)
			return
		}
		filename = existingWardrobe.Filename
	} else {
		http.Error(w, "Не удалось получить файл", http.StatusBadRequest)
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

	wardrobe.ID = id

	err = Service.UpdateWardrobe(&wardrobe, fileBytes)
	if err != nil {
		switch err {
		case Model.ErrWardrobeNotFound:
			http.Error(w, "Шкаф не найден", http.StatusNotFound)
		case Model.ErrInvalidWardrobeData:
			http.Error(w, "Неверные данные шкафа", http.StatusBadRequest)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": wardrobe.Title + " успешно отредактирован!",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteWardrobeHandler(w http.ResponseWriter, r *http.Request) {
	if !Service.CheckTokenExpiration(w, r) {
		return
	}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Неверный ID шкафа", http.StatusBadRequest)
		return
	}

	err = Service.DeleteWardrobe(id)
	if err != nil {
		switch err {
		case Model.ErrWardrobeNotFound:
			http.Error(w, "Шкаф не найден", http.StatusNotFound)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Шкаф успешно удален!",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetWardrobeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Неверный ID шкафа", http.StatusBadRequest)
		return
	}

	wardrobe, err := Service.GetWardrobeById(id)
	if err != nil {
		switch err {
		case Model.ErrWardrobeNotFound:
			http.Error(w, "Шкаф не найден", http.StatusNotFound)
		default:
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wardrobe)
}
