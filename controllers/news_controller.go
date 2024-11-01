package controllers

import (
    "github.com/gorilla/mux"
    "encoding/json"
    "net/http"
    "github.com/imrezaulkrm/bartadhara/database"
    "strconv"
)

// সমস্ত নিউজ ফেচ করার জন্য ফাংশন
func GetAllNews(w http.ResponseWriter, r *http.Request) {
    news, err := database.FetchAllNews()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(news)
}

// নির্দিষ্ট নিউজ ফেচ করার জন্য ফাংশন
func GetNewsByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    news, err := database.FetchNewsByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(news)
}

// নতুন নিউজ তৈরি করার জন্য ফাংশন
func CreateNews(w http.ResponseWriter, r *http.Request) {
    var news database.News
    if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := database.InsertNews(news); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
}

// নিউজ আপডেট করার জন্য ফাংশন
func UpdateNews(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var news database.News
    if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := database.UpdateNews(id, news); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

// নিউজ ডিলিট করার জন্য ফাংশন
func DeleteNews(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    if err := database.DeleteNews(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
