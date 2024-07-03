package Router

import (
	"github.com/gorilla/mux"
	"tests/Controller"
)

func AuthRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/registration", Controller.Registration).Methods("POST")
	r.HandleFunc("/login", Controller.Login).Methods("POST")
	return r
}
