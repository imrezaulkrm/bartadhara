package routes

import (
    "github.com/gorilla/mux"
    "news_backend/controllers"
)

func UserRoutes(r *mux.Router) {
    r.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
    r.HandleFunc("/login", controllers.LoginUser).Methods("POST")
}
