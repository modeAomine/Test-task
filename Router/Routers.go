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

func AdminUserRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/user/update/{id}", Middleware.AuthMiddleware(http.HandlerFunc(Controller.UpdateUser))).Methods("PUT")
	r.Handle("/user/delete/{id}", Middleware.AuthMiddleware(http.HandlerFunc(Controller.DeleteUser))).Methods("DELETE")

	return r
}

func MixRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/auth").Handler(AuthRouter())
	r.PathPrefix("/wardrobe").Handler(WardrobeRouter())
	r.PathPrefix("/user").Handler(AdminUserRouter())
	return r
}
