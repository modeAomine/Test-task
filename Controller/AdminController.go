package Controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tests/Model"
	"tests/Service"
)

func UpdateUserByAdmin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updateRequest struct {
		User Model.User `json:"user"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := updateRequest.User
	user.ID = id

	err = Service.UpdateUserByAdmin(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"message": "User updated successfully", "user": user}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
