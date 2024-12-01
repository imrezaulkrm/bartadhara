package routes

import (
    "github.com/gorilla/mux"
    "github.com/imrezaulkrm/bartadhara-backend/controllers"
)

// AdminRoutes defines admin-related routes
func AdminRoutes(r *mux.Router) {
    ac := controllers.AdminController{}

    // Admin related routes
    r.HandleFunc("/admin", ac.FetchAllAdmins).Methods("GET") // Get all admins
    r.HandleFunc("/admin/{id}", ac.FetchAdminByID).Methods("GET") // Get an admin by ID
    //r.HandleFunc("/admin", ac.InsertAdmin).Methods("POST") // Admin registration route
    r.HandleFunc("/admin/{id}", ac.UpdateAdmin).Methods("PUT") // Admin update route
    r.HandleFunc("/admin/{id}", ac.DeleteAdmin).Methods("DELETE") // Admin delete route

    // Admin authentication routes
    r.HandleFunc("/admin/login", ac.Login).Methods("POST") // Admin login route
    r.HandleFunc("/admin/register", ac.Register).Methods("POST") // Admin registration route
}
