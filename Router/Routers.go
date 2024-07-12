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

func UserRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user/update/{id}", Controller.UpdateUserProfile).Methods("PUT")
	r.Handle("/user/all", Middleware.AuthMiddleware(http.HandlerFunc(Controller.GetAllUsers))).Methods("GET")

	return r
}

func AdminWardrobeRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/admin/wardrobe/add", Middleware.AuthMiddleware(http.HandlerFunc(Controller.AddWardrobeHandler))).Methods("POST")
	r.Handle("/admin/wardrobe/update/{id}", Middleware.AuthMiddleware(http.HandlerFunc(Controller.UpdateWardrobeHandler))).Methods("PUT")
	r.Handle("/admin/wardrobe/delete/{id}", Middleware.AuthMiddleware(http.HandlerFunc(Controller.DeleteWardrobeHandler))).Methods("DELETE")
	r.Handle("/admin/wardrobe/all", Middleware.AuthMiddleware(http.HandlerFunc(Controller.GetAllWardrobe))).Methods("GET")

	return r
}

func AdminUserRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/admin/user/update/{id}", Middleware.AuthMiddleware(http.HandlerFunc(Controller.UpdateUserByAdmin))).Methods("PUT")
	r.Handle("/admin/user/delete/{id}", Middleware.AuthMiddleware(http.HandlerFunc(Controller.DeleteUser))).Methods("DELETE")
	r.Handle("/admin/user/add", Middleware.AuthMiddleware(http.HandlerFunc(Controller.CreateUser))).Methods("POST")

	return r
}

func MixRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/auth").Handler(AuthRouter())
	r.PathPrefix("/admin").Handler(AdminWardrobeRouter())
	r.PathPrefix("/admin").Handler(AdminUserRouter())
	r.PathPrefix("/user").Handler(UserRouter())

	return r
}
