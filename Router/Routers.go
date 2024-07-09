package Router

import (
	"github.com/gorilla/mux"
	"net/http"
	"tests/Controller"
	"tests/Middleware"
)

func AuthRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/auth/registration", Controller.Registration).Methods("POST")
	r.HandleFunc("/auth/login", Controller.Login).Methods("POST")
	return r
}

func WardrobeRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/wardrobe", Middleware.AuthMiddleware(http.HandlerFunc(Controller.AddWardrobeHandler))).Methods("POST")

	return r
}

func MixRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/auth").Handler(AuthRouter())
	r.PathPrefix("/wardrobe").Handler(WardrobeRouter())
	return r
}

func wardrobeHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*Middleware.Claims)
	w.Write([]byte("Hello, " + claims.Role))
}
