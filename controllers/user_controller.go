package controllers

import (
    "encoding/json"
    "net/http"
    "news_backend/models"
    "news_backend/database"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    database.DB.Create(&user)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
    // লগইন লজিক এখানে লেখো
}
