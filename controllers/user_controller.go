package controllers

import (
    "log"
    "github.com/gorilla/mux"
    "encoding/json"
    "net/http"
    "strconv"
    //"database/sql"

    "news_backend/database"
    "news_backend/models"
)

// UserController struct
type UserController struct{}

// FetchAllUsers handles GET requests to fetch all users
func (uc *UserController) FetchAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := database.FetchAllUsers()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(users)
}

// FetchUserByID handles GET requests to fetch a user by ID
func (uc *UserController) FetchUserByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil || id <= 0 {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    log.Printf("Fetching user with ID: %d", id) // লগিং যুক্ত করুন

    user, err := database.FetchUserByID(id)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound) // সঠিক স্ট্যাটাস কোড
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}


// InsertUser handles POST requests to create a new user
func (uc *UserController) InsertUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    dbUser := database.User{
        Username: user.Username,
        Email:    user.Email,
        Password: user.Password,
        Picture:  user.Picture,
    }

    err = database.InsertUser(dbUser)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// UpdateUser handles PUT requests to update a user
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var user models.User
    err = json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    dbUser := database.User{
        Username: user.Username,
        Email:    user.Email,
        Password: user.Password,
        Picture:  user.Picture,
    }

    err = database.UpdateUser(id, dbUser) // Ensure this function exists in your database.go
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// DeleteUser handles DELETE requests to remove a user
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    err = database.DeleteUser(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// নির্দিষ্ট ইউজার ফেচ করার জন্য ফাংশন
func GetUserByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil || id <= 0 {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    user, err := database.FetchUserByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}