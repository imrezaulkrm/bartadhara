package routes

import (
    "github.com/gorilla/mux"
    "news_backend/controllers"
)

// UserRoutes defines user-related routes
func UserRoutes(r *mux.Router) {
    uc := controllers.UserController{}

    r.HandleFunc("/users", uc.FetchAllUsers).Methods("GET")
    r.HandleFunc("/users/{id}", uc.FetchUserByID).Methods("GET")
    r.HandleFunc("/users", uc.InsertUser).Methods("POST")
    r.HandleFunc("/users/{id}", uc.UpdateUser).Methods("PUT")
    r.HandleFunc("/users/{id}", uc.DeleteUser).Methods("DELETE")
}
