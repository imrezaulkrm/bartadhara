package routes

import (
    "github.com/gorilla/mux"
    "github.com/imrezaulkrm/bartadhara/controllers"
)

// UserRoutes defines user-related routes
func UserRoutes(r *mux.Router) {
    uc := controllers.UserController{}

    // User related routes
    r.HandleFunc("/users", uc.FetchAllUsers).Methods("GET")
    r.HandleFunc("/users/{id}", uc.FetchUserByID).Methods("GET")
    r.HandleFunc("/users", uc.InsertUser).Methods("POST") // Registration route
    r.HandleFunc("/users/{id}", uc.UpdateUser).Methods("PUT")
    r.HandleFunc("/users/{id}", uc.DeleteUser).Methods("DELETE")

    // Authentication routes
    r.HandleFunc("/login", uc.Login).Methods("POST") // Login route
    // New routes for user categories
    r.HandleFunc("/users/{id}/categories", uc.FetchUserCategories).Methods("GET")
    r.HandleFunc("/users/{id}/categories", uc.UpdateUserCategories).Methods("PUT")
}

